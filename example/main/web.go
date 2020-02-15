package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	// "github.com/opentracing-contrib/go-gin/ginhttp"
	
	// ...

	"github.com/hlf513/go-micro-pkg/config"
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
	config.Init(etcdAddress)

	etcd := etcdv3.NewRegistry(func(op *registry.Options) {op.Addrs = config2.GetEtcd().Host})

	service := web.NewService(
		web.Name(config2.GetServer().Name),
		web.Version("latest"),
		web.Registry(etcd),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal("service.Init:", err.Error())
	}

	// create router
	engine := gin.New()
	engine.Use(gin.Recovery())
	// jaeger 
	t, closer := tracer.Connect(config2.GetJaeger().Name, config2.GetJaeger().Address)
	defer closer.Close()
	engine.Use(ginhttp.Middleware(t))
	engine.Use(gin2.Metadata())
	// zap 
	engine.Use(gin2.Ginzap())
	engine.Use(gin2.RecoveryWithZap())

	router.Register(engine)

	// register handler
	service.Handle("/", engine)
	// run service
	if err := service.Run(); err != nil {
		log.Fatal("service.Run", err.Error())
	}
}
