package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/database/mysql"
	"go-filestore-server/service/account/handler"
	proto "go-filestore-server/service/account/proto"
	dbproxy "go-filestore-server/service/dbproxy/client"
	"log"
	"time"
)

func init() {
	config.InitConfig("./service/bin/config.json")
	mysql.InitMysql()
}

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.DefaultConfig.ConsulAddr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	// 初始化service, 解析命令行参数等
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
