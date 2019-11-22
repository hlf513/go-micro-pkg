package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type trace struct {
	span opentracing.Span
}

func Trace(ctx context.Context) *trace {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return &trace{span: span}
	}

	return &trace{}
}

func (t trace) Log(logs ...interface{}) {
	if t.span != nil {
		t.span.LogKV(logs...)
	}
}

func (t trace) Tag(key string, value interface{}) {
	if t.span != nil {
		t.span.SetTag(key, value)
	}
}

func (t trace) Sampler() {
	if t.span != nil{
		ext.SamplingPriority.Set(t.span, 1)
	}	
}
