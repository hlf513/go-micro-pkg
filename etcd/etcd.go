package etcd

import (
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"

	"github.com/hlf513/go-micro-pkg/config"
)

// Init 初始化 etcd
func Init() registry.Registry {
	return etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.GetEtcd().Host,
		}
		op.Timeout = config.GetEtcd().Timeout * time.Second
	})
}
