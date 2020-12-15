package model

import (
	"fmt"
	"go-filestore-server/database/mysql"
	"go-filestore-server/logger"
	"time"
)

// UserFile: 用户文件表结构体
type UserFile struct {
	UserName   string `json:"UserName"`   // 用户名
	FileHash   string `json:"FileHash"`   // 文件Hash
	FileName   string `json:"FileName"`   // 文件名
	FileSize   int64  `json:"FileSize"`   // 文件尺寸
	UploadAt   string `json:"UploadAt"`   // 上传时间
	LastUpdate string `json:"LastUpdate"` // 最后一次更新时间
}

// 添加: 插入用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	insertSQL := "insert ignore into tbl_user_file (`user_name`,`file_hash`,`file_name`,`file_size`,`upload_at`) values(?,?,?,?,?) "
	// insert ignore 会忽略数据库中已经存在的数据
	fmt.Println(insertSQL)
	stmt, err := mysql.DBConn().Prepare(insertSQL)
	if err != nil {
		fmt.Println("prepare err:\t", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now())
	if err != nil {
		fmt.Println("exec err:\t", err.Error())
		return false
	}
	return true
}

// 查询: 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	querySQL := "select file_hash, file_name, file_size, upload_at,last_update from tbl_user_file where user_name=? limit ?"
	stmt, err := mysql.DBConn().Prepare(querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdate)
		if err != nil {
			logger.Info(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}

// 删除文件
func DeleteUserFile(username, filehash string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set status=2 where user_name=? and file_hash=? limit 1")
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash)
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	return true
}

// 更新: 文件重命名
func RenameFileName(username, filehash, filename string) bool {
	stmt, err := mysql.DBConn().Prepare("update tbl_user_file set file_name = ? where user_name = ? and file_hash = ? limit 1")
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash)
	if err != nil {
		logger.Info(err.Error())
		return false
	}
	return true
}

// 查询: 用户单个文件信息
func QueryUserFileMeta(username string, filehash string) (*UserFile, error) {
	stmt, err := mysql.DBConn().Prepare("select file_hash, file_name, file_size, upload_at, last_update from tbl_user_file where user_name = ? and file_hash =? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		return nil, err
	}

	ufile := UserFile{}
	if rows.Next() {
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdate)
		if err != nil {
			logger.Info(err.Error())
			return nil, err
		}
	}
	return &ufile, nil
}
