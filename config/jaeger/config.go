package jaeger

import "github.com/hlf513/go-micro-pkg/config"

type Jaeger struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// j 初始化
var j = &Jaeger{}

// GetJaegerConf 读取配置
func GetJaegerConf() (*Jaeger, error) {
	if err := config.GetConfigurator().Get([]string{"jaeger"}, j); err != nil {
		return nil, err
	}
	return j, nil
}
