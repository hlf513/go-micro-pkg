package redis

import (
	"time"

	"github.com/hlf513/go-micro-pkg/config"
)

type Redis struct {
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

// rs 初始化
var rs = make(map[string]Redis)

// GetRedis 读取配置
func GetRedisConf() (map[string]Redis, error) {
	err := config.GetConfigurator().Get([]string{"redis"}, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
