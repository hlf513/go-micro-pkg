module go-micro-pkg

go 1.13

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/getsentry/sentry-go v0.3.1
	github.com/gin-gonic/gin v1.4.0
	github.com/gogo/protobuf v1.3.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.11
	github.com/micro/go-micro v1.16.0
	github.com/micro/go-plugins v1.5.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.uber.org/zap v1.13.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
