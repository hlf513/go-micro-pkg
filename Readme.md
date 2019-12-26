# go-micro-pkg

## 组件列表
1. Log => zap
2. trace => jaeger
3. etcd
4. mysql
5. redis
6. sentry

## 中间件
1. web 
2. srv

## 使用方法

### 系统配置
####  配置文件
todo
#### micro.web
todo
#### micro.srv
todo
### 组件使用

#### Redis
```go
// 从连接池获取连接
rds := redis.GetConn(c.Request.Context())
// 连接放回连接池
defer rds.Close()
```

#### MySQL
```go
// 从连接池获取 DB 连接
db := mysql.GetDB(ctx)
```

### Sentry
```go
// 只打印文字信息
log.SentryMessage(message string)
// 打印文字信息 + 调用堆栈
log.SentryException(err error)
```

### Zap
```go
// warn,debug,fatal,error
log.Info(ctx,"信息内容","分类")
```