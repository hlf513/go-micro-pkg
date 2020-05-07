package etcd

import (
	"time"

	"github.com/micro/go-micro/v2/registry"
	e "github.com/micro/go-micro/v2/registry/etcd"
)

// DefaultRegistry 注册 etcd
func DefaultRegistry() registry.Registry {
	conf := GetConf()
	return e.NewRegistry(func(op *registry.Options) {
		op.Addrs = conf.Host
		op.Timeout = conf.Timeout * time.Second
	})
}
