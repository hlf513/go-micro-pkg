package ginTrace

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
)

// GinInjectMetadata 注入 metadata 到 tracing（原生 micro srv 的 trace 依赖 metadata）
func InjectMetadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		md := make(map[string]string)
		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
		if err != nil {
			log.Fatal("[jaeger metadata]:", err.Error())
		}
		// ctx := context.TODO()
		ctx := c.Request.Context()
		ctx = metadata.NewContext(ctx, md)
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
