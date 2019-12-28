package zap

import (
	"encoding/json"
	"sync"

	"github.com/hlf513/go-micro-pkg/config"
)

// zapLogger 定义配置项
type zapLogger struct {
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool
	Level      string
}

// conf 定义更新配置
type conf struct {
	ZapLogger *zapLogger
}

var (
	zapConf = &conf{}
	s       sync.RWMutex
)

// GetConf 读取配置
func GetConf() *zapLogger {
	s.RLock()
	defer s.RUnlock()

	return zapConf.ZapLogger
}

// SetConf 更新配置
func SetConf(c []byte) error {
	s.Lock()
	defer s.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"zap"}, &zapConf.ZapLogger); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, zapConf); err != nil {
			return err
		}
	}

	return nil
}
