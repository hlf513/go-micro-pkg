package main

import (
	"os"
	"strings"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/service/handler"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"

)

var etcdAddress []string

func init() {
	etcdAddr := os.Getenv("micro_config_etcd_address")
	if etcdAddr == "" {
		log.Fatal("请设置环境变量：micro_config_etcd_address")
	}
	etcdAddress = strings.Split(etcdAddr, ",")
	log.Info("etcd_address:", etcdAddress)
}

func main() {
	// 配置初始化
	conf.Init(etcdAddress)

	// 使用 etcd 注册
	etcd := etcdv3.NewRegistry(func(op *registry.Options) { op.Addrs = config.GetEtcd().Host })

	t, closer := tracer.Connect(config.GetJaeger().Name, config.GetJaeger().Address)
	defer closer.Close()
	// 初始化 service
	service := micro.NewService(
		micro.Name(config.GetServer().Name),
		micro.Registry(etcd),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		micro.WrapHandler(micro2.NewTraceWrapper()),
		
		micro.WrapHandler(micro2.WaitGroupWrapper(micro2.WaitGroup())),
		micro.AfterStop(func() error {
			micro2.WaitGroup().Wait()
			time.Sleep(1 * time.Second)
			return nil
		}),
	)
	service.Init()

	// 注册 Handler
	handler.Register(service)

	// 运行 service
	if err := service.Run(); err != nil {
		log.Fatal("server run:", err.Error())
	}
}
