package jaeger

import (
	"context"

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
