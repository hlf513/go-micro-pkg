package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
)

// Metadata 为 micro.web 增加 metadata 
func Metadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		md := make(map[string]string)
		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
		if err != nil {
			log.Info("jaeger.metadata:", err.Error())
		}
		ctx := context.TODO()
		ctx = metadata.NewContext(ctx, md)
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if c.Writer.Status() > 200 {
			span := opentracing.SpanFromContext(ctx)
			if span != nil {
				span.SetTag("error", true)
			}
		}
	}
}
