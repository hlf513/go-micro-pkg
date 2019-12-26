package mysql

import (
	"time"

	"github.com/hlf513/go-micro-pkg/config"
)

type DB struct {
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

// dbs 初始化
var dbs = make(map[string]DB)

// GetDBs 读取配置
func GetDBConf() (map[string]DB, error) {
	err := config.GetConfigurator().Get([]string{"db"}, &dbs)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}
