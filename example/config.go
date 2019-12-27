package example

import (
	"path/filepath"

	"github.com/micro/go-micro/config/source/etcd"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"

	util "github.com/hlf513/go-micro-pkg/config"
	"github.com/hlf513/go-micro-pkg/config/mysql"
	"github.com/hlf513/go-micro-pkg/config/redis"
	"github.com/hlf513/go-micro-pkg/config/zap"
)

// Init 配置文件初始化
func Init() {
	// 文件源
	var sp = string(filepath.Separator)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))
	configPath := filepath.Join(appPath, "/config/config.yaml")
	fileSource := file.NewSource(
		file.WithPath(configPath),
	)
	// etcd 源
	etcdSource := etcd.NewSource(
		etcd.WithAddress("127.0.0.1:2379"),
		etcd.WithPrefix("/micro/config/aidc"),
		etcd.StripPrefix(true), // 返回值过滤前缀
	)
	// file 会覆盖 etcd
	if err := util.GetConfigurator().Init(etcdSource, fileSource); err != nil {
		log.Fatal(err.Error())
	}

	// 日志初始化
	zap.Setup()
	// MySQL
	mysql.Connect()
	// Redis
	redis.Connect()
	// ...
}
