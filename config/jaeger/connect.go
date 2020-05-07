package jaeger

import (
	"io"

	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
	client "github.com/uber/jaeger-client-go"
)

// Connect 创建一个jaeger Tracer
func Connect() (opentracing.Tracer, io.Closer) {
	conf := GetConf()
	sender, err := client.NewUDPTransport(conf.Address, 0)
	if err != nil {
		log.Fatal("[jaeger] connect was failed ", err.Error())
	}

	var tracer, closer = client.NewTracer(
		conf.Name,
		GetSampler(conf.Env),
		client.NewRemoteReporter(sender),
	)

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}

func Close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Error("[jaeger] close was failed ", err.Error())
	}
}
