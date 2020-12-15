package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/database/ceph"
	"go-filestore-server/database/mq"
	"go-filestore-server/database/oss"
	"go-filestore-server/logger"
	"go-filestore-server/meta"
	"go-filestore-server/model"
	"go-filestore-server/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Gin版本
// 响应上传页面
func UploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/upload.html")
	if err != nil {
		c.String(404, `网页不存在`)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

// 处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "upload failed",
			})
		}
	}()

	// 1.从form表单中获得文件内容句柄
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		logger.Info("failed to get form data, err:\t", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	// 2.把文件内容转为[]byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		logger.Info("failed to get file data,err:\t", err.Error())
		errCode = -2
		return
	}

	// 3.构建文件元信息
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		FileHash: util.Sha1(buf.Bytes()),
		FileSize: int64(len(buf.Bytes())),
		UploadAt: time.Now().Format(common.StandardTimeFormat),
	}

	// 4.将文件写入临时存储位置
	fileMeta.Location = config.DefaultConfig.TempLocalRootDir + fileMeta.FileHash
	fmt.Println("file meta:\t", fileMeta, head.Size)

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Println("failed to create file, err:\t", err.Error())
		errCode = -3
		return
	}
	defer newFile.Close()

	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != fileMeta.FileSize || err != nil {
		fmt.Println("failed to save data into file, writtenSize:\t", nByte, "err:\t", err.Error())
		errCode = -4
		return
	}

	// 5.同步或异步将文件转移到Ceph/OSS
	newFile.Seek(0, 0) // 游标重新回到文件头部
	if config.DefaultConfig.CurrentStoreType == common.StoreCeph {
		fmt.Println("走ceph")
		// 文件写入Ceph存储
		data, _ := ioutil.ReadAll(newFile)
		cephPath := "/ceph/" + fileMeta.FileHash
		_ = ceph.PutObject("userfile", cephPath, data)
		fileMeta.Location = cephPath
	} else if config.DefaultConfig.CurrentStoreType == common.StoreOSS {
		fmt.Println("走oss")
		// 文件写入OSS存储
		ossPath := "oss/" + fileMeta.FileHash
		// 判断写入OSS为同步还是异步
		if !config.DefaultConfig.AsyncTransferEnable {
			// TODO 设置oss的文件名， 方便指定文件名下载
			err = oss.Bucket().PutObject(ossPath, newFile)
			if err != nil {
				logger.Info(err.Error())
				errCode = -5
				return
			}
			fileMeta.Location = ossPath
		} else {
			// 写入异步转移任务队列
			data := mq.TransferData{
				FileHash:      fileMeta.FileHash,
				CurLocation:   fileMeta.Location,
				DestStoreType: common.StoreOSS,
				DestLocation:  ossPath,
			}
			pubData, _ := json.Marshal(data)
			pubSuc := mq.Publish(config.DefaultConfig.TransExchangeName,
				config.DefaultConfig.TransOSSRoutingKey,
				pubData)
			if !pubSuc {
				// TODO 当前发送转移信息失败，稍后重试
			}
		}
	}

	// 6.更新文件表记录
	_ = meta.UpdateFileMetaDB(fileMeta)
	fmt.Println("更新用户文件表")
	// 7.更新用户文件表
	username := c.Request.FormValue("username")
	suc := model.OnUserFileUploadFinished(username, fileMeta.FileHash, fileMeta.FileName, fileMeta.FileSize)
	fmt.Println(suc)
	if suc {
		c.Redirect(http.StatusFound, "/static/view/home.html")
	} else {
		errCode = -6
	}
}

// 上传已经完成
func UploadSucHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "upload finish",
	})
	return
}

// 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -2,
			"msg":  "upload failed!",
		})
		return
	}

	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": -3,
				"msg":  "upload failed!",
			})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -4,
			"msg":  "no such file",
		})
	}
}

// 批量查询文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	if limitCnt == 0 {
		limitCnt = 10
	}
	userFiles, err := model.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "query failed!",
			"err":  err.Error(),
		})
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -2,
			"msg":  "query failed!",
		})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

// 文件下载接口
func DownloadHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	// TODO: 处理异常情况
	fm, _ := meta.GetFileMetaDB(filehash)
	userFile, _ := model.QueryUserFileMeta(username, filehash)

	if strings.HasPrefix(fm.Location, config.DefaultConfig.TempLocalRootDir) {
		c.FileAttachment(fm.Location, userFile.FileName)
	} else if strings.HasPrefix(fm.Location, config.DefaultConfig.CephRootDir) {
		// ceph中的文件，通过ceph api先下载
		bucket := ceph.GetCephBucket("userfile")
		data, _ := bucket.Get(fm.Location)
		c.Header("content-disposition", "attachment;filename=\""+userFile.FileName+"\"")
		c.Data(http.StatusOK, "application/octect-stream", data)
	}
}

// 更新元信息接口
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	filehash := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	newFileName := c.Request.FormValue("filename")

	if opType != "0" || len(newFileName) < 1 {
		c.Status(http.StatusForbidden)
		return
	}

	// 更新用户文件表tbl_user_file中的文件名，tbl_file的文件名不用修改
	_ = model.RenameFileName(username, filehash, newFileName)

	// 返回最新的文件信息
	userFile, err := model.QueryUserFileMeta(username, filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFile)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

// 删除文件及元信息
func FileDeleteHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	fm, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// 删除本地文件
	os.Remove(fm.Location)
	// TODO: 可以考虑删除Ceph/OSS上的文件
	// 可以不立即删除，加个超时机制
	// 比如该文件10天后也没有用户再次上传，那么就可以真正的删除了

	// 删除文件表中的一条记录
	suc := model.DeleteUserFile(username, filehash)
	if !suc {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// 尝试秒传接口
func TryFastUploadHandler(c *gin.Context) {
	// 1.解析请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 2.从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		logger.Info(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3.查不到记录则返回秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	// 4.上传过则将文件信息写入到用户表，返回成功
	suc := model.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	return
}

// 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	// 从文件表查找记录
	row, _ := model.GetFileMeta(filehash)

	// TODO 判断文件存储在OSS，还是Ceph，还是在本地
	if strings.HasPrefix(row.FileAddr.String, config.DefaultConfig.TempLocalRootDir) ||
		strings.HasPrefix(row.FileAddr.String, config.DefaultConfig.CephRootDir) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		tmpURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Data(http.StatusOK, "octet-stream", []byte(tmpURL))
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// oss 下载url
		signedURL := oss.DownloadURL(row.FileAddr.String)
		logger.Info(row.FileAddr.String)
		c.Data(http.StatusOK, "octet-stream", []byte(signedURL))
	}
}
