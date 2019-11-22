package config

import "time"

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

var dbs = make(map[string]DB)

func GetDBs() map[string]DB {
	return dbs
}

func SetDBs(s string, d DB) {
	dbs[s] = d
}
