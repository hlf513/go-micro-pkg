package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	"github.com/opentracing-contrib/go-gin/ginhttp"

	e "github.com/hlf513/go-micro-pkg/config/etcd"
	"github.com/hlf513/go-micro-pkg/config/jaeger"
	"github.com/hlf513/go-micro-pkg/config/server"
	"github.com/hlf513/go-micro-pkg/example/config"
	"github.com/hlf513/go-micro-pkg/example/main/web/router"
	"github.com/hlf513/go-micro-pkg/wrapper/common"
	"github.com/hlf513/go-micro-pkg/wrapper/web/ginTrace"
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

	service := web.NewService(
		web.Name(server.GetConf().Name),
		web.Version("latest"),

		web.Registry(etcd.NewRegistry(func(op *registry.Options) {
			op.Addrs = e.GetConf().Host
		})),

		// web.RegisterTTL(30*time.Second),
		// web.RegisterInterval(15*time.Second),

		web.BeforeStop(func() error {
			common.WaitGroup().Wait()
			time.Sleep(1 * time.Second) // 等待响应输出后再关闭
			return nil
		}),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal("service.Init:", err.Error())
	}

	// Jaeger
	t, closer := jaeger.Connect()
	defer jaeger.Close(closer)

	// create a gin engine
	engine := gin.New()
	// trace
	engine.Use(ginhttp.Middleware(t))
	engine.Use(ginTrace.InjectMetadataWrapper())
	engine.Use(ginTrace.LogToTraceWrapper())
	engine.Use(ginTrace.RecoveryWrapper())

	router.Register(engine)

	// register handler
	service.Handle("/", engine)
	// run service
	if err := service.Run(); err != nil {
		log.Fatal("service.Run", err.Error())
	}
}
