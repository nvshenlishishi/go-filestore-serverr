package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ratelimit2 "github.com/juju/ratelimit"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
	"go-filestore-server/common"
	"go-filestore-server/config"
	userProto "go-filestore-server/service/account/proto"
	dlProto "go-filestore-server/service/download/proto"
	upProto "go-filestore-server/service/upload/proto"
	"go-filestore-server/util"
	"log"
	"net/http"
	"strconv"
)

var (
	userCli userProto.UserService
	upCli   upProto.UploadService
	dlCli   dlProto.DownloadService
)

func init() {
	config.InitConfig("./service/bin/config.json")
	bRate := ratelimit2.NewBucketWithRate(100, 1000)
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.DefaultConfig.ConsulAddr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.apigw"),
		micro.Flags(common.CustomFlags...),
		micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)), // 加入限流功能, false为不等待(超限即返回请求失败)
		micro.WrapClient(hystrix.NewClientWrapper()),               // 加入熔断功能, 处理rpc调用失败的情况(cirucuit breaker)
	)

	service.Init()

	cli := service.Client()

	userCli = userProto.NewUserService("go.micro.service.user", cli)
	upCli = upProto.NewUploadService("go.micro.service.upload", cli)
	dlCli = dlProto.NewDownloadService("go.micro.service.download", cli)
}

func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: passwd,
	})

	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}

func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func DoSigninHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	resp, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: username,
		Password: password,
	})

	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if resp.Code != common.StatusOK {
		c.JSON(200, gin.H{
			"msg":  "登录失败",
			"code": resp.Code,
		})
		return
	}

	cliResp := util.RespMsg{
		Code: int(common.StatusOK),
		Msg:  "登陆成功",
		Data: struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location:      "/static/view/home.html",
			Username:      username,
			Token:         resp.Token,
			UploadEntry:   config.DefaultConfig.UploadMicroHost,
			DownloadEntry: config.DefaultConfig.DownloadServiceHost,
		},
	}

	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}

func UserInfoHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	resp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	fmt.Println("用户信息:\t", resp.SignupAt, resp.Username, resp.Status, resp.Message, resp.Code)
	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"user_name":   username,
			"signup_at":   resp.SignupAt,
			"last_active": resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")

	rpcResp, err := userCli.UserFiles(context.TODO(), &userProto.ReqUserFile{
		Username: username,
		Limit:    int32(limitCnt),
	})
	fmt.Println(string(rpcResp.FileData))

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}
	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	fileSha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	newFileName := c.Request.FormValue("filename")

	if opType != "0" || len(newFileName) < 1 {
		c.Status(http.StatusForbidden)
		return
	}

	rpcResp, err := userCli.UserFileRename(context.TODO(), &userProto.ReqUserFileRename{
		Username:    username,
		Filehash:    fileSha1,
		NewFileName: newFileName,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}
	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}
