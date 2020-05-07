package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	s "github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"

	"github.com/hlf513/go-micro-pkg/config/etcd"
	"github.com/hlf513/go-micro-pkg/config/jaeger"
	"github.com/hlf513/go-micro-pkg/config/server"
	"github.com/hlf513/go-micro-pkg/example/config"
	"github.com/hlf513/go-micro-pkg/example/main/srv/handler"
	"github.com/hlf513/go-micro-pkg/wrapper/common"
	"github.com/hlf513/go-micro-pkg/wrapper/srv"
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
	config.Init(etcdAddress)

	// Jaeger
	t, closer := jaeger.Connect()
	defer jaeger.Close(closer)

	// 初始化 service
	service := micro.NewService(
		micro.Name(server.GetConf().Name),
		micro.Registry(etcdv3.NewRegistry(func(op *registry.Options) {
			op.Addrs = etcd.GetConf().Host
		})),

		micro.Version("latest"),
		// micro.RegisterTTL(30*time.Second),
		// micro.RegisterInterval(15*time.Second),

		// trace
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		micro.WrapHandler(srv.TraceLogWrapper()),

		// recovery
		micro.WrapHandler(srv.RecoveryWrapper()),

		// waitGroup && graceful shutdown
		micro.WrapHandler(srv.WaitGroupWrapper(common.WaitGroup())),
		micro.BeforeStop(func() error {
			common.WaitGroup().Wait()
			time.Sleep(1 * time.Second) // 等待响应输出后再关闭
			return nil
		}),
	)
	service.Init()

	// 注册 Handler
	handler.Register(service)

	// register subscriber
	// topic 不能和 server name 相同；否则会重复消费
	micro.RegisterSubscriber(
		"test",
		service.Server(),
		new(Sub),
		s.SubscriberQueue("test-queue"), // 增加这个参数后，多个 sub 只有一个会执行
	)

	// 运行 service
	if err := service.Run(); err != nil {
		log.Fatal("server run:", err.Error())
	}
}

// All methods of Sub will be executed when
// a message is received
type Sub struct{}

// Method can be of any name
// common.Event 自定义的 proto
func (s *Sub) Process(ctx context.Context, event *common.Event) error {
	// md, _ := metadata.FromContext(ctx)
	log.Logf("[pubsub.1] Received event %+v with metadata\n", event)
	// do something with event
	return nil
}
