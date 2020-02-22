package config

import (
	"path/filepath"

	"github.com/micro/go-micro/config/source/etcd"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"

	util "github.com/hlf513/go-micro-pkg/config"
	"github.com/hlf513/go-micro-pkg/config/hystrix"
	"github.com/hlf513/go-micro-pkg/config/jaeger"
	"github.com/hlf513/go-micro-pkg/config/mysql"
	"github.com/hlf513/go-micro-pkg/config/redis"
	"github.com/hlf513/go-micro-pkg/config/sentry"
	"github.com/hlf513/go-micro-pkg/config/server"
	"github.com/hlf513/go-micro-pkg/config/zap"
	"github.com/hlf513/go-micro-pkg/utils"
)

// Init 配置文件初始化
func Init(etcdAddress []string) {
	// 文件源
	appPath := utils.CurrentDir()
	configPath := filepath.Join(appPath, "/config/example.yaml")
	fileSource := file.NewSource(
		file.WithPath(configPath),
	)
	// etcd 源
	etcdSource := etcd.NewSource(
		etcd.WithAddress(etcdAddress...),
		etcd.WithPrefix("/micro/config/aidc"),
		etcd.StripPrefix(true), // 返回值过滤前缀
	)
	// file 会覆盖 etcd
	if err := util.GetConfigurator().Init(
		[]util.SetterFunc{
			mysql.SetConf,
			jaeger.SetConf,
			redis.SetConf,
			sentry.SetConf,
			server.SetConf,
			zap.SetConf,
			hystrix.SetConf,
		},
		etcdSource,
		fileSource,
	); err != nil {
		log.Fatal(err.Error())
	}

	// Zap
	// zap.Setup()
	// MySQL
	// mysql.Connect()
	// Redis
	// redis.Connect()
	// ...
}
