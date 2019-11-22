package zap

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go-micro-pkg/config"
)

var zlogger *zap.Logger

func MineTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func Setup(l config.Logger) {
	config.SetLogger(l)
	hook := lumberjack.Logger{
		Filename:   l.FilePath,   // 日志文件路径
		MaxSize:    l.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: l.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     l.MaxAge,     // 文件最多保存多少天
		Compress:   l.Compress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "datetime",
		LevelKey:       "lvl",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     MineTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch l.Level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	default:
		atomicLevel.SetLevel(zap.ErrorLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),               // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), //  打印到文件
		atomicLevel,                                         // 日志级别
	)

	zlogger = zap.New(core)
}

func Close() {
	_ = zlogger.Sync()
}

func GetLogger() *zap.Logger {
	return zlogger
}
