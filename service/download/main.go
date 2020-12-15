package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go-filestore-server/common"
	"go-filestore-server/config"
	dbproxy "go-filestore-server/service/dbproxy/client"
	dlProto "go-filestore-server/service/download/proto"
	"go-filestore-server/service/download/route"
	dlRpc "go-filestore-server/service/download/rpc"
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
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(config.DefaultConfig.DownloadServiceHost)
}

func init() {
	config.InitConfig("./service/bin/config.json")
}
func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
