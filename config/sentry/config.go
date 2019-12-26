package sentry

import "github.com/hlf513/go-micro-pkg/config"

type Sentry struct {
	Dns string `json:"dns"`
}

// sentry 初始化
var ss = &Sentry{}

// GetSentryConf 读取配置
func GetSentryConf(s Sentry) (*Sentry, error) {
	if err := config.GetConfigurator().Get([]string{"sentry"}, ss); err != nil {
		return nil, err
	}
	return ss, nil
}
