# github.com/hlf513/go-micro-pkg

此项目主要用于 go-micro 技术栈，作为框架的扩展存在。

## 组件列表
1. Log => zap
2. trace => jaeger
3. etcd
4. mysql
5. redis
6. sentry

## 使用方法

### 系统配置
####  配置文件
> config/config.yaml
```yaml
server:
  name: go.micro.srv.entity
  env: dev
  mode: srv
db:
  default:
    type: mysql
    host: 127.0.0.1:3306
    username: test
    password: test
    dbname: example
    max_idle_conn: 10
    max_open_conn: 10
    max_lifetime: 30
    # true 输出 SQL
    debug: true
  slave:
    type: mysql
    host: 127.0.0.1:3306
    username: test
    password: test
    dbname: example
    max_idle_conn: 10
    max_open_conn: 10
    max_lifetime: 30
    debug: true
redis:
  default:
    server: 127.0.0.1:6379
    max_idle: 1
    max_active: 50
    connect_timeout: 500
    read_timeout: 300
    write_timeout: 300
    idle_timeout: 240
    password:
    select_db: 0
etcd:
  host: http://127.0.0.1:2379
log:
  file_path: /var/log/example.log
  # 单位 M
  max_size: 1  
  max_backups: 30
  max_age: 7
  compress: false
  # debug | info | warn | error
  level: debug
sentry:
  dns: http://22cceca506e44660ac58c6a29fa7e583@mac.mini:9999/2
jaeger:
  name: trace-name
  address: localhost:6831
```

> config/config.go
```go
package config

import (
	"path/filepath"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/util/log"

	c "github.com/hlf513/go-micro-pkg/config"
	"github.com/hlf513/go-micro-pkg/mysql"
	"github.com/hlf513/go-micro-pkg/sentry"
	"github.com/hlf513/go-micro-pkg/zap"

	"github.com/hlf513/go-micro-pkg/redis"
)

type configs struct {
	Server c.Server
	DBs    map[string]c.DB `json:"db"`
	Redis  map[string]c.Redis
	Etcd   c.Etcd
	Log    c.Logger
	Sentry c.Sentry
	Jaeger c.Jaeger
}

// conf 配置信息（不能并发读写）
var conf configs

// Init 配置文件初始化
func Init() {
	var sp = string(filepath.Separator)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("."+sp, sp)))
	configPath := filepath.Join(appPath, "/config/config.yaml")

	// 加载配置文件
	if err := config.LoadFile(configPath); err != nil {
		log.Fatal("未找到配置文件:"+configPath+";err:"+err.Error(), "application")
	}

	// 解析配置文件
	if err := config.Scan(&conf); err != nil {
		log.Fatal("解析配置文件异常:"+err.Error(), "application")
	}

	c.SetServer(conf.Server)
	c.SetJaeger(conf.Jaeger)

	// 日志初始化
	zap.Setup(conf.Log)
	// db
	mysql.Connect(conf.DBs)
	// redis
	redis.Connect(conf.Redis)
	// sentry
	sentry.Init(conf.Sentry)
}
```
#### micro.web
```go
package main

import (
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	openTrace "github.com/opentracing/opentracing-go"

	config2 "github.com/hlf513/go-micro-pkg/config"
	gin2 "github.com/hlf513/go-micro-pkg/middleware/gin"
	"github.com/hlf513/go-micro-pkg/etcd"
	"github.com/hlf513/go-micro-pkg/jaeger"
	// config
	// router
)

func main() {
	// 配置初始化
	config.Init()

	// 服务初始化
	service := web.NewService(
		web.Name(config2.GetServer().Name),
		web.Version("latest"),
		web.Registry(etcd.Init()),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
	)
	if err := service.Init(); err != nil {
		log.Fatal("service.Init:", err.Error())
	}
	
	// jaeger
	t, closer := jaeger.Connect(config2.GetJaeger().Name, config2.GetJaeger().Address)
	defer closer.Close()

	// 路由初始化
	engine := gin.New()
	engine.Use(gin.Recovery())
	// sentry middleware
	engine.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	// trace middleware
	engine.Use(ginhttp.Middleware(t))
	engine.Use(gin2.Metadata())
	// zap middleware
	engine.Use(gin2.Ginzap())
	engine.Use(gin2.RecoveryWithZap())
	// 
	router.Register(engine)

	// register handler
	service.Handle("/", engine)
	
	// run service
	if err := service.Run(); err != nil {
		log.Fatal("service.Run", err.Error())
	}
}
```

#### micro.srv
```go
package main

import (
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	openTrace "github.com/opentracing/opentracing-go"

	config2 "github.com/hlf513/go-micro-pkg/config"
	tracer "github.com/hlf513/go-micro-pkg/jaeger"
	micro2 "github.com/hlf513/go-micro-pkg/middleware/micro"
	"github.com/hlf513/go-micro-pkg/etcd"
	// config 
	// register 
)


func main() {
	// 配置初始化
	config.Init()

	// jaeger
	t, closer := jaeger.Connect(config2.GetJaeger().Name, config2.GetJaeger().Address)
	defer closer.Close()
	// 初始化 service
	service := micro.NewService(
		micro.Name(config2.GetServer().Name),
		micro.Registry(etcd.Init()),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapHandler(opentracing.NewHandlerWrapper(t)),
		micro.WrapHandler(micro2.NewTraceWrapper()),
	)
	service.Init()

	// 注册 Handler
	register.Init(service)

	// 运行 service
	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
```

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
fatal、warn、error 会自动上报 sentry
```go
// warn,debug,fatal,error
log.Info(ctx,"信息内容","分类")
```

