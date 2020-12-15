package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/database/mq"
	"go-filestore-server/database/oss"
	dbproxy "go-filestore-server/service/dbproxy/client"
	"log"
	"os"
	"time"
)

func init() {
	config.InitConfig("./service/bin/config.json")
}

func main() {
	// 文件转移服务
	go startTranserService()

	// rpc 服务
	startRPCService()
}

func startTranserService() {
	if !config.DefaultConfig.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")

	// 初始化mq client
	mq.InitMq()

	mq.StartConsume(
		config.DefaultConfig.TransOSSQueueName,
		"transfer_oss",
		Transfer)
}

func startRPCService() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.DefaultConfig.ConsulAddr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.transfer"), // 服务名称
		micro.RegisterTTL(time.Second*10),       // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),   // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}),
	)

	// 初始化dbproxy client
	dbproxy.Init(service)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// Transfer : 处理文件转移
func Transfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	resp, err := dbproxy.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !resp.Suc {
		log.Println("更新数据库异常，请检查:" + pubData.FileHash)
		return false
	}
	return true
}
