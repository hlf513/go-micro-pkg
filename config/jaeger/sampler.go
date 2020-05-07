package jaeger

import (
	client "github.com/uber/jaeger-client-go"

	envc "github.com/hlf513/go-micro-pkg/const/env"
)

// sampler 采样率
var sampler client.Sampler

// SetSampler 设置采样率
func SetSampler(s client.Sampler) {
	sampler = s
}

// GetSampler 获取采样率
func GetSampler(env string) client.Sampler {
	if env == envc.Development {
		return client.NewConstSampler(true) // 全量追踪
	} else if sampler == nil {
		SetProbabilisticSample(0.01) // 1%的采样率（报错时100%记录）
	}

	return sampler
}

// SetProbabilisticSample 设置采样率（0.0~1.0)
func SetProbabilisticSample(samplingRate float64) {
	sampler, _ = client.NewProbabilisticSampler(samplingRate)
}
