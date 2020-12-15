package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/logger"
	"go-filestore-server/model"
	"go-filestore-server/util"
	"net/http"
	"time"
)

// Gin
// 响应注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// 处理注册Post请求
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "请求参数无效",
			"code": common.StatusParamInvalid,
		})
		return
	}

	// 对密码进行加盐及取sha1值加密
	encPasswd := util.Sha1([]byte(passwd + config.DefaultConfig.PwdSalt))
	// 将用户信息注册到用户表
	suc := model.UserSignup(username, encPasswd)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册成功",
			"code": common.StatusOK,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册失败",
			"code": common.StatusRegisterFailed,
		})
	}
}

// 响应登陆页面
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// 处理登陆post请求
func DoSignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	encPasswd := util.Sha1([]byte(password + config.DefaultConfig.PwdSalt))
	logger.Info("signIn:\t", username, password, encPasswd)
	// 1.校验用户名及密码
	pwdChecked := model.UserSignin(username, encPasswd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "登陆失败",
			"code": common.StatusLoginFailed,
		})
		return
	}

	// 2.生成访问凭证token
	token := GenToken(username)
	upRes := model.UpdateToken(username, token)
	if !upRes {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "登陆失败",
			"code": common.StatusLoginFailed,
		})
		return
	}

	// 3.登陆成功后重定向到首页
	resp := util.RespMsg{
		Code: int(common.StatusOK),
		Msg:  "登陆成功",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	//	token := c.Request.FormValue("token")

	// 2. 查询用户信息
	user, err := model.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	// 3.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.Data(http.StatusOK, "octet-stream", resp.JSONBytes())
}

// 查询用户是否存在
func UserExistHandler(c *gin.Context) {
	// 1.解析请求参数
	username := c.Request.FormValue("username")

	// 2.查询用户信息
	exists, err := model.UserExist(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.StatusServerError,
			"msg":  "server error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":   common.StatusOK,
			"msg":    "ok",
			"exists": exists,
		})
	}
}

// 生成Token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
