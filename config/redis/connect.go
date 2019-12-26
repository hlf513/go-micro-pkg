package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/micro/go-micro/util/log"
)

// rds 连接池
var rds = make(map[string]*redis.Pool)

// Connect 创建 Redis 连接池
func Connect() {
	configs, err := GetRedisConf()
	if err != nil {
		log.Fatal("[redis connect] ", err.Error())
	}

	for name, config := range configs {
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
		log.Info("[redis connect] 初始化 Reids 连接：" + name)
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

// getPool 获取 Redis 连接池
func getPool(name string) *redis.Pool {
	if r, ok := rds[name]; ok {
		return r
	} else {
		log.Fatal("[redis GetPool] 未找到 redis 连接池:" + name)
	}
	return nil
}
