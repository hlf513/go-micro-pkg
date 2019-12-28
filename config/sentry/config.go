package sentry

import (
	"encoding/json"
	"sync"

	"github.com/hlf513/go-micro-pkg/config"
)

// sentry 定义配置项
type sentry struct {
	Dns string `json:"dns"`
}

// conf 定义更新配置
type conf struct {
	Sentry *sentry
}

var (
	sentryConf = &conf{}
	s          sync.RWMutex
)

// GetConf 读取配置
func GetConf() *sentry {
	s.RLock()
	defer s.RUnlock()

	return sentryConf.Sentry
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"sentry"}, &sentryConf.Sentry); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, sentryConf); err != nil {
			return err
		}
	}

	return nil
}
