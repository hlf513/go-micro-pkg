package micro

import (
	"context"
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"
)

// NewTraceWrapper 为 micro.rpc 增加请求响应日志
func NewTraceWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			span := opentracing.SpanFromContext(ctx)
			if span != nil {
				request, _ := json.Marshal(req.Body())
				span.LogKV("request", string(request))
				opentracing.ContextWithSpan(ctx, span)
			}

			err := h(ctx, req, rsp)

			if span != nil {
				if err != nil {
					span.SetTag("error", true)
					span.LogKV("error_msg", err.Error())
				}

				m := jsonpb.Marshaler{
					EmitDefaults: true,
					OrigName:     true,
					EnumsAsInts:  true,
				}

				response, _ := m.MarshalToString(rsp.(proto.Message))
				span.LogKV("response", string(response))
				opentracing.ContextWithSpan(ctx, span)
			}

			return err
		}
	}
}
