package redis

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hlf513/go-micro-pkg/config"
)

// rds 定义配置项
type rds struct {
	Server         string
	MaxIdle        int           `json:"max_idle"`
	MaxActive      int           `json:"max_active"`
	ConnectTimeout time.Duration `json:"connect_timeout"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	IdleTimeout    time.Duration `json:"idle_timeout"`
	Password       string        `json:"password"`
	SelectDB       int           `json:"select_db"`
}

// conf 定义更新配置
type conf struct {
	Redis map[string]rds `json:"redis"`
}

var (
	redisConf = &conf{}
	s         sync.RWMutex
)

// GetConf 读取配置
func GetConf() map[string]rds {
	s.RLock()
	defer s.RUnlock()

	return redisConf.Redis
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"redis"}, &redisConf.Redis); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, redisConf); err != nil {
			return err
		}
	}

	return nil
}
