package srv

import (
	"context"
	e "errors"
	"fmt"
	"runtime/debug"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"

	"github.com/hlf513/go-micro-pkg/config/jaeger"
)

// RecoveryWrapper 恢复 srv 的 panic 异常，并返回 599 异常
func RecoveryWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
			defer func() {
				if p := recover(); p != nil {
					span := opentracing.SpanFromContext(ctx)
					errMsg := fmt.Sprintf("%v", p)
					if span != nil {
						jaeger.SetError(span, e.New(errMsg))
						// 记录错误日志
						span.LogKV(
							"error",
							fmt.Sprintf(
								`[Recovery from panic] - %s - %s`,
								errMsg,
								debug.Stack(),
							))
						opentracing.ContextWithSpan(ctx, span)
					}
					err = errors.New(req.Method(), errMsg, 599)
				}
			}()

			err = h(ctx, req, rsp)

			return err
		}
	}
}
