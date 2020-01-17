package hystrix

import (
	"encoding/json"
	"sync"

	"github.com/hlf513/go-micro-pkg/config"
)

type hystrix struct {
	// 是否打开 Streaming event 服务
	StreamServer bool `json:"stream_server"`
	// Streaming event 服务端口
	ServerPort string `json:"server_port"`
	// 重试次数
	Retries int `json:"retries"`
	// 超时时间
	Timeout int `json:"timeout"`
	// 最大并发数
	MaxConcurrent int `json:"max_concurrent"`
	// 判断前至少请求数
	VolumeThreshold int `json:"volume_threshold"`
	// 熔断后睡眠时间
	SleepWindow int `json:"sleep_window"`
	// 错误请求百分比
	ErrorPercentThreshold int `json:"error_percent_threshold"`
}

// conf 定义更新配置
type conf struct {
	Hystrix *hystrix
}

var (
	hystrixConf = &conf{}
	hs          sync.RWMutex
)

// GetConf 读取配置
func GetConf() *hystrix {
	hs.RLock()
	defer hs.RUnlock()

	return hystrixConf.Hystrix
}

// SetConf 更新配置
func SetConf(c []byte) error {
	hs.Lock()
	defer hs.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"hystrix"}, &hystrixConf.Hystrix); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, hystrixConf); err != nil {
			return err
		}
	}

	return nil
}
