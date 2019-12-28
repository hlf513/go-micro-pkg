package jaeger

import client "github.com/uber/jaeger-client-go"

// sampler 采样率
var sampler client.Sampler

// SetSampler 设置采样率
func SetSampler(s client.Sampler) {
	sampler = s
}

// GetSampler 获取采样率
func GetSampler() client.Sampler {
	if sampler == nil {
		return client.NewConstSampler(true) // 全量追踪
	}

	return sampler
}
