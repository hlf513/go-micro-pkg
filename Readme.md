# go-micro-pkg

## Dependent
1. Log => zap
2. trace => jaeger
3. etcd
4. mysql
5. redis
6. sentry

## Trace Middleware
1. web 
2. srv

## Util
1. 分布式锁 - [etcd](https://github.com/Scalingo/go-etcd-lock)
2. Http lib 
3. AES (ECB 加密模式)

## 使用方法

### 动态配置
```bash
# /example/config.go
```

### Redis
```go
// 从连接池获取连接
rds := redis.GetConn(c.Request.Context())
// 连接放回连接池
defer rds.Close()
```

### MySQL
```go
// 从连接池获取 DB 连接
db := mysql.GetDB(ctx)
```

### Sentry
```go
// 只打印文字信息
sentry.SentryMessage(message string)
// 打印文字信息 + 调用堆栈
sentry.SentryException(err error)
```

### Zap
```go
// 增加字段（注意不会替换已存在字段）
zap.SetOtherFields([]zap.Field{})
// warn,debug,fatal,error
zap.Info(ctx,"信息内容")
```