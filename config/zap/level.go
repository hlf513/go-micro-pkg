package zap

import (
	"context"

	"go.uber.org/zap"

	"github.com/hlf513/go-micro-pkg/config/jaeger"
	"github.com/hlf513/go-micro-pkg/config/server"
)

// fields 存储 message 字段
var fields []zap.Field

// SetOtherFields 设置 fields
func SetOtherFields(fs []zap.Field) {
	fields = append(fields, fs...)
}

// otherFields 获取其他字段
func otherFields(ctx context.Context) []zap.Field {
	serverConf := server.GetConf()
	otherFields := []zap.Field{
		zap.String("traceId", jaeger.GetTraceId(ctx)),
		zap.String("app", serverConf.Name),
		zap.String("env", serverConf.Env),
	}

	if len(fields) > 0 {
		otherFields = append(otherFields, fields...)
	}

	return otherFields
}

// Info 打印 info 日志
func Info(ctx context.Context, msg string) {
	GetLogger().Info(msg, otherFields(ctx)...)
}

// Warn 打印 warn 日志
func Warn(ctx context.Context, msg string) {
	GetLogger().Warn(msg, otherFields(ctx)...)
}

// Error 打印 error 日志
func Error(ctx context.Context, msg string) {
	GetLogger().Error(msg, otherFields(ctx)...)
}

// Debug 打印 debug 日志
func Debug(ctx context.Context, msg string) {
	GetLogger().Debug(msg, otherFields(ctx)...)
}

// Fatal 调用此方法 Fatal 级别会退出，Panic 级别会执行 panic()
func Fatal(ctx context.Context, msg string) {
	GetLogger().Fatal(msg, otherFields(ctx)...)
}
