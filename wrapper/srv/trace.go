package srv

import (
	"context"
	"encoding/json"
	"runtime/debug"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// TraceLog 记录 rpc server 的请求和响应到 tracing 
func TraceLog() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			span := opentracing.SpanFromContext(ctx)
			if span != nil {
				// 记录请求
				request, _ := json.Marshal(req.Body())
				span.LogKV("request", string(request))
				opentracing.ContextWithSpan(ctx, span)
			}

			err := h(ctx, req, rsp)

			if span != nil {
				if err != nil {
					// 记录错误信息
					span.SetTag("error", true)
					ext.SamplingPriority.Set(span, 1)
					span.LogKV("error_msg", err.Error(), "debug.stack", debug.Stack())
				}

				// 记录响应
				m := jsonpb.Marshaler{
					EmitDefaults: true,
					OrigName:     true,
					EnumsAsInts:  true,
				}
				response, _ := m.MarshalToString(rsp.(proto.Message))
				span.LogKV("response", response)
				opentracing.ContextWithSpan(ctx, span)
			}

			return err
		}
	}
}
