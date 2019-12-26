package etcd

import (
	"time"

	"github.com/hlf513/go-micro-pkg/config"
)

type Etcd struct {
	Host    []string
	Timeout time.Duration
}

// etcd 初始化
var etcd = &Etcd{}

// GetEtcdConf 读取配置
func GetEtcdConf() (*Etcd, error) {
	if err := config.GetConfigurator().Get([]string{"etcd"}, etcd); err != nil {
		return nil, err
	}
	return etcd, nil
}
