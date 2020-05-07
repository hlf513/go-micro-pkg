package jaeger

import (
	"context"

	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	client "github.com/uber/jaeger-client-go"
)

// GetTraceId 获取 trace id
func GetTraceId(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return span.Context().(client.SpanContext).TraceID().String()
	}

	return ""
}

// StartTrace 开启请求追踪
func StartTrace(operationName string) (opentracing.Span, context.Context) {
	span, c := opentracing.StartSpanFromContext(context.Background(), operationName)
	md := make(map[string]string)
	_ = opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	ctx := opentracing.ContextWithSpan(metadata.NewContext(c, md), span)

	return span, ctx
}

// StopTrace 停止请求追踪
func StopTrace(span opentracing.Span) {
	if span != nil {
		span.Finish()
	}
}

// SetError 设置异常请求标签
func SetError(span opentracing.Span) {
	if span != nil {
		span.SetTag("error", true)
	}
}
