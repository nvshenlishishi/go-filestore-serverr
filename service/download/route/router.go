package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-filestore-server/service/download/api"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://localhost:8080"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))
	// 文件下载相关接口
	router.GET("/file/download", api.DownloadHandler)
	router.POST("/file/downloadurl", api.DownloadURLHandler)

	return router
}
