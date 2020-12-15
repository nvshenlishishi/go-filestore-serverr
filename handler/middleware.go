package handler

import (
	"github.com/gin-gonic/gin"
	"go-filestore-server/common"
	"go-filestore-server/util"
	"net/http"
)

// gin版本
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		// 验证登陆token是否有效
		if len(username) < 3 || !IsTokenValid(token) {
			c.Abort()
			resp := util.NewRespMsg(int(common.StatusTokenInvalid), "token无效", nil)
			c.JSON(http.StatusOK, resp)
			// c.Redirect(http.StatusFound, "/static/view/signin.html")
			return
		}
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,Range,X-Requested-with, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
