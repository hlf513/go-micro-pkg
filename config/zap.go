package config

type Logger struct {
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool
	Level      string
}

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func GetLogger() Logger {
	return logger
}
