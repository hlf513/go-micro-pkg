package zap

import "github.com/hlf513/go-micro-pkg/config"

type ZapLogger struct {
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool
	Level      string
}

// zapLogger 初始化
var zapLogger = &ZapLogger{}

// GetLoggerConf 读取配置
func GetLoggerConf() (*ZapLogger, error) {
	if err := config.GetConfigurator().Get([]string{"zap"}, zapLogger); err != nil {
		return nil, err
	}
	return zapLogger, nil
}
