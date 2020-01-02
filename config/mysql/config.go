package mysql

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hlf513/go-micro-pkg/config"
)

// DB 定义配置项
type db struct {
	Type        string
	Host        string
	Username    string
	Password    string
	DBName      string
	MaxIdleConn int           `json:"max_idle_conn"`
	MaxOpenConn int           `json:"max_open_conn"`
	MaxLifeTime time.Duration `json:"max_lifetime"`
	Debug       bool
}

// conf 定义更新配置
type conf struct {
	DB map[string]db `json:"db"`
}

var (
	dbConf = &conf{}
	s      sync.RWMutex
)

// GetConf 读取配置
func GetConf() map[string]db {
	s.RLock()
	defer s.RUnlock()

	return dbConf.DB
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"db"}, &dbConf.DB); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, dbConf); err != nil {
			return err
		}
	}

	return nil
}
