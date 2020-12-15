package main

import (
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/database/mysql"
	dbProxy "go-filestore-server/service/dbproxy/proto"
	dbRpc "go-filestore-server/service/dbproxy/rpc"
	"log"
	"time"
)

func init() {
	config.InitConfig("./service/bin/config.json")
	mysql.InitMysql()
}

func main() {
	startRpcService()
}

func startRpcService() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.DefaultConfig.ConsulAddr,
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.dbproxy"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),      // 声明超时时间, 避免consul不主动删掉已失去心跳的服务节点
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init(micro.Action(func(c *cli.Context) {
		dbHost := c.String("dbhost")
		if len(dbHost) > 0 {
			log.Println("custom db address: " + dbHost)
			// UpdateDBHost(dbHost)
		}
	}))

	mysql.InitMysql()

	dbProxy.RegisterDBProxyServiceHandler(service.Server(), new(dbRpc.DBProxy))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

var MySQLSource = "test:test@tcp(127.0.0.1:3306)/fileserver?charset=utf8"

func UpdateDBHost(host string) {
	MySQLSource = fmt.Sprintf("test:test@tcp(%s)/fileserver?charset=utf8", host)
}
