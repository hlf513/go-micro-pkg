package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// GetConn 获取 Redis 连接
func GetConn(ctx context.Context, poolName ...string) *redisConn {
	var p string
	if len(poolName) == 0 {
		p = "default"
	} else {
		p = poolName[0]
	}

	conn := getPool(p).Get()

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

// redisConn 
type redisConn struct {
	conn redis.Conn
	span opentracing.Span
}

// Do 执行 redis 命令
func (r redisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	if r.span != nil {
		r.span.LogKV("cmd", commandName, "args", args)
	}
	replay, err := r.conn.Do(commandName, args...)
	if err != nil {
		if r.span != nil {
			ext.SamplingPriority.Set(r.span, 1)
			r.span.SetTag("error", true)
		}
		return nil, err
	}
	if r.span != nil {
		r.span.LogKV("value", replay)
		r.span.Finish()
	}
	return replay, nil
}

// Close 关闭 redis 连接
func (r redisConn) Close() {
	if r.span != nil {
		r.span.Finish()
	}
	_ = r.conn.Close()
}
