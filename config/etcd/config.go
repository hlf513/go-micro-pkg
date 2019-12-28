package etcd

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hlf513/go-micro-pkg/config"
)

// etcd 定义配置项
type etcd struct {
	Host    []string
	Timeout time.Duration
}

// conf 定义更新配置
type conf struct {
	Etcd *etcd
}

var (
	etcdConf = &conf{}
	s        sync.RWMutex
)

// GetConf 读取配置
func GetConf() *etcd {
	s.RLock()
	defer s.RUnlock()

	return etcdConf.Etcd
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"etcd"}, &etcdConf.Etcd); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, etcdConf); err != nil {
			return err
		}
	}

	return nil
}
