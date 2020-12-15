package main

import (
	"fmt"
	"go-filestore-server/config"
	"go-filestore-server/database/mysql"
	"go-filestore-server/database/redisPool"
	"go-filestore-server/logger"
	"go-filestore-server/route"
	"os"
)

func main() {
	config.InitConfig("./config/config.json")
	InitTemp()
	logger.Init()
	mysql.InitMysql()
	redisPool.InitRedis()
	// gin framework
	router := route.Router()

	// 启动服务并监听端口
	err := router.Run(config.DefaultConfig.UploadServiceHost)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}

func InitTemp() {
	// 目录已存在
	if _, err := os.Stat(config.DefaultConfig.TempLocalRootDir); err == nil {
		fmt.Println("目录已经存在:\t", config.DefaultConfig.TempLocalRootDir)
		return
	} else {
		fmt.Println("目录不存在:\t", err.Error())
	}

	// 尝试创建目录
	err := os.MkdirAll(config.DefaultConfig.TempLocalRootDir, 0744)
	if err != nil {
		fmt.Println("无法创建临时存储目录，程序将退出")
		os.Exit(1)
	}
}
