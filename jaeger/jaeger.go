package jaeger

import (
	"context"
	"io"

	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)


var sampler jaeger.Sampler

// Connect 创建一个jaeger Tracer
func Connect(serverName, host string) (opentracing.Tracer, io.Closer) {
	sender, err := jaeger.NewUDPTransport(host, 0)
	if err != nil {
		log.Fatal("jaeger.connect was failed")
	}
	var tracer, closer = jaeger.NewTracer(
		serverName,
		GetSampler(),
		jaeger.NewRemoteReporter(sender),
	)
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func SetSampler(s jaeger.Sampler) {
	sampler = s
}

func GetSampler() jaeger.Sampler {
	if sampler == nil {
		return jaeger.NewConstSampler(true) // 全量追踪
	}

	return sampler
}

func GetTraceId(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return span.Context().(jaeger.SpanContext).TraceID().String()
	}

	return ""
}