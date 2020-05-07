package jaeger

import (
	"encoding/json"
	"sync"

	"github.com/hlf513/go-micro-pkg/config"
)

// jaeger 定义配置项
type jaeger struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Env     string `json:"env"`
}

// conf 定义更新配置
type conf struct {
	Jaeger *jaeger
}

var (
	jaegerConf = &conf{}
	s          sync.RWMutex
)

// GetConf 读取配置
func GetConf() *jaeger {
	s.RLock()
	defer s.RUnlock()

	return jaegerConf.Jaeger
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"jaeger"}, &jaegerConf.Jaeger); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, jaegerConf); err != nil {
			return err
		}
	}

	return nil
}
