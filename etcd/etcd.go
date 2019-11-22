package etcd

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"

	"go-micro-pkg/config"
)

// Init 初始化 etcd
func Init() registry.Registry {
	return etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.GetEtcd().Host,
		}
	})
}
