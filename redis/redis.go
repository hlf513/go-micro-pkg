package redis

import (
	"time"

	"github.com/micro/go-micro/util/log"

	conf "go-micro-pkg/config"

	"github.com/gomodule/redigo/redis"
)

// rds 只读
var rds = make(map[string]*redis.Pool)

// Connect 创建 Redis 连接池
func Connect(configs map[string]conf.Redis) {
	for name, config := range configs {
		conf.SetRedis(name, config)
		redisPool := &redis.Pool{
			MaxIdle:     config.MaxIdle,
			MaxActive:   config.MaxActive,
			IdleTimeout: config.IdleTimeout * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(
					"tcp",
					config.Server,
					redis.DialConnectTimeout(config.ConnectTimeout*time.Millisecond),
					redis.DialReadTimeout(config.ReadTimeout*time.Millisecond),
					redis.DialWriteTimeout(config.WriteTimeout*time.Millisecond),
				)
				if err != nil {
					return nil, err
				}
				if config.Password != "" {
					if _, err := c.Do("AUTH", config.Password); err != nil {
						_ = c.Close()
						return nil, err
					}
				}
				if config.SelectDB != 0 {
					if _, err := c.Do("SELECT", config.SelectDB); err != nil {
						_ = c.Close()
						return nil, err
					}
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			Wait: true,
		}

		rds[name] = redisPool
		log.Info("初始化 Reids 连接：" + name)
	}
}

// Close 关闭 Redis 连接池
func Close() {
	for _, r := range rds {
		if r != nil {
			_ = r.Close()
		}
	}
}

// GetPool 从 Redis 连接池获取 Redis 连接
func GetPool(name string) *redis.Pool {
	if r, ok := rds[name]; ok {
		return r
	} else {
		log.Fatal("未找到 redis 连接池:"+name, "application")
	}
	return nil
}
