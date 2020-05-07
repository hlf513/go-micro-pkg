package hystrix

import (
	hg "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
	micro "github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
)

// DefaultRetries 默认重试次数
const DefaultRetries = 3

// initConf 根据参数初始化 hystrix
func initConf() {
	conf := GetConf()
	if conf.Timeout > 1000 {
		hg.DefaultTimeout = conf.Timeout
	}
	if conf.MaxConcurrent > 0 {
		hg.DefaultMaxConcurrent = conf.MaxConcurrent
	}
	if conf.VolumeThreshold > 0 {
		hg.DefaultVolumeThreshold = conf.VolumeThreshold
	}
	if conf.SleepWindow > 1000 {
		hg.DefaultSleepWindow = conf.SleepWindow
	}
	if conf.ErrorPercentThreshold > 0 {
		hg.DefaultErrorPercentThreshold = conf.ErrorPercentThreshold
	}
}

// Client 获取 hystrix client
func Client() client.Client {
	conf := GetConf()
	if conf == nil {
		return client.DefaultClient
	}

	initConf()

	var retries = conf.Retries
	if retries == 0 {
		retries = DefaultRetries
	}

	cli := micro.NewClientWrapper()(client.DefaultClient)
	_ = cli.Init(client.Retries(retries))

	return cli
}
