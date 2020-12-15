package main

import (
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/database/mq"
	dbproxy "go-filestore-server/service/dbproxy/client"
	upProto "go-filestore-server/service/upload/proto"
	"go-filestore-server/service/upload/route"
	upRpc "go-filestore-server/service/upload/rpc"
	"log"
	"os"
	"time"
)

func startRPCService() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.DefaultConfig.ConsulAddr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
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
	// 初始化mq client
	mq.InitMq()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	fmt.Println("upload addr:\t", config.DefaultConfig.UploadMicroHost)
	router.Run(config.DefaultConfig.UploadMicroHost)
	// service := web.NewService(
	// 	web.Name("go.micro.web.upload"),
	// 	web.Handler(router),
	// 	web.RegisterTTL(10*time.Second),
	// 	web.RegisterInterval(5*time.Second),
	// )
	// if err := service.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := service.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

func init() {
	config.InitConfig("./service/bin/config.json")
	os.MkdirAll(config.DefaultConfig.TempLocalRootDir, 0777)
	os.MkdirAll(config.DefaultConfig.TempPartRootDir, 0777)
}

func main() {
	// api 服务
	go startAPIService()
	// rpc 服务
	startRPCService()
}
