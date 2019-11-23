package log

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/hlf513/go-micro-pkg/config"
	"github.com/hlf513/go-micro-pkg/jaeger"
	zap2 "github.com/hlf513/go-micro-pkg/zap"
)

var fields []zap.Field

func SetFields(fds []zap.Field) {
	fields = fds
}

func otherFields(ctx context.Context, category string) []zap.Field {
	if len(fields) == 0 {
		return []zap.Field{
			zap.String("traceId", jaeger.GetTraceId(ctx)),
			zap.String("app", config.GetServer().Name),
			zap.String("env", config.GetServer().Env),
			zap.String("mode", config.GetServer().Mode),
			zap.String("category", category),
		}
	} else {
		return fields
	}
}

func Info(ctx context.Context, msg, category string) {
	zap2.GetLogger().Info(msg, otherFields(ctx, category)...)
}

func Warn(ctx context.Context, msg, category string) {
	SentryException(errors.New(msg))
	zap2.GetLogger().Warn(msg, otherFields(ctx, category)...)
}

func Error(ctx context.Context, msg, category string) {
	SentryException(errors.New(msg))
	zap2.GetLogger().Error(msg, otherFields(ctx, category)...)
}

func Debug(ctx context.Context, msg, category string) {
	zap2.GetLogger().Debug(msg, otherFields(ctx, category)...)
}

func Fatal(ctx context.Context, msg, category string) {
	SentryException(errors.New(msg))
	zap2.GetLogger().Fatal(msg, otherFields(ctx, category)...)
}
