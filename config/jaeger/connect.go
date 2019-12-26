package jaeger

import (
	"io"

	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

// Connect 创建一个jaeger Tracer
func Connect() (opentracing.Tracer, io.Closer) {
	conf, err := GetJaegerConf()
	if err != nil {
		log.Fatal("[jaeger connect]", err.Error())
	}

	sender, err := jaeger.NewUDPTransport(conf.Address, 0)
	if err != nil {
		log.Fatal("[jaeger connect] connect was failed", err.Error())
	}

	var tracer, closer = jaeger.NewTracer(
		conf.Name,
		GetSampler(),
		jaeger.NewRemoteReporter(sender),
	)

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}
