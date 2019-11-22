package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/opentracing/opentracing-go"
)

// GetConn 获取 Redis 连接
func GetConn(ctx context.Context, poolName ...string) *redisConn {
	var p string
	if len(poolName) == 0 {
		p = "default"
	} else {
		p = poolName[0]
	}

	conn := GetPool(p).Get()

	var span opentracing.Span
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		span = parentSpan.Tracer().StartSpan("Redis", opentracing.ChildOf(parentSpan.Context()))
	}

	return &redisConn{
		conn: conn,
		span: span,
	}
}

type redisConn struct {
	conn redis.Conn
	span opentracing.Span
}

func (r redisConn) Close() {
	if r.span != nil {
		r.span.Finish()
	}
	_ = r.conn.Close()
}

func (r redisConn) Set(key string, value interface{}) error {
	r.span.SetTag("command", "set")
	r.span.SetTag("key", key)
	r.span.SetTag("value", value)
	if _, err := r.conn.Do("Set", key, value); err != nil {
		r.span.SetTag("error", true)
		return err
	}
	return nil
}

func (r redisConn) SetExpire(key string, value interface{}, expire int) error {
	r.span.SetTag("command", "setExpire")
	r.span.SetTag("key", key)
	r.span.SetTag("value", value)
	r.span.SetTag("expire", expire)
	if _, err := r.conn.Do("Set", key, value, "EX", expire); err != nil {
		r.span.SetTag("error", true)
		return err
	}
	return nil
}

func (r redisConn) GetString(key string) (string, error) {
	r.span.SetTag("command", "getString")
	r.span.SetTag("key", key)
	rep, err := redis.String(r.conn.Do("Get", key))
	if err == redis.ErrNil {
		r.span.SetTag("value", "")
		return "", nil
	}
	if err != nil {
		r.span.SetTag("error", true)
		return "", err
	}
	r.span.SetTag("value", rep)
	return rep, nil
}

func (r redisConn) GetBytes(key string) ([]byte, error) {
	r.span.SetTag("command", "GetBytes")
	r.span.SetTag("key", key)
	rep, err := redis.Bytes(r.conn.Do("Get", key))
	if err == redis.ErrNil {
		r.span.SetTag("value", nil)
		return nil, nil
	}
	if err != nil {
		r.span.SetTag("error", true)
		return nil, err
	}
	r.span.SetTag("value", rep)
	return rep, nil
}

func (r redisConn) Del(key string) error {
	r.span.SetTag("command", "Del")
	r.span.SetTag("key", key)
	if _, err := r.conn.Do("Del", key); err != nil {
		r.span.SetTag("error", true)
		return err
	}
	return nil
}

func (r redisConn) Lock(key string, timeout int) (bool, error) {
	r.span.SetTag("command", "Lock")
	r.span.SetTag("key", key)
	r.span.SetTag("timeout", key)
	reply, err := r.conn.Do("Set", key, 1, "EX", timeout, "NX")
	if err != nil {
		r.span.SetTag("error", true)
		return false, err
	}
	r.span.SetTag("value", reply.(bool))
	return reply.(bool), nil
}

func (r redisConn) Unlock(key string) error {
	r.span.SetTag("command", "UnLock")
	return r.Del(key)
}
