package config

import "time"

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

var redis = make(map[string]Redis)

func GetRedis() map[string]Redis {
	return redis
}

func SetRedis(s string, r Redis) {
	redis[s] = r
}
