package etcd

import (
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

// DefaultRegistry 注册 etcd
func DefaultRegistry() registry.Registry {
	conf := GetConf()
	return etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = conf.Host
		op.Timeout = conf.Timeout * time.Second
	})
}
