package server

import (
	"encoding/json"
	"sync"

	"github.com/hlf513/go-micro-pkg/config"
)

// server 定义配置项
type server struct {
	Name string
	Env  string
}

// conf 定义更新配置
type conf struct {
	Server *server
}

var (
	serverConf = &conf{}
	s          sync.RWMutex
)

// GetConf 读取配置
func GetConf() *server {
	s.RLock()
	defer s.RUnlock()

	return serverConf.Server
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"server"}, &serverConf.Server); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, serverConf); err != nil {
			return err
		}
	}

	return nil
}
