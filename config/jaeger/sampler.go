package jaeger

import "github.com/uber/jaeger-client-go"

// sampler 采样率
var sampler jaeger.Sampler

// SetSampler 设置采样率
func SetSampler(s jaeger.Sampler) {
	sampler = s
}

// GetSampler 获取采样率
func GetSampler() jaeger.Sampler {
	if sampler == nil {
		return jaeger.NewConstSampler(true) // 全量追踪
	}

	return sampler
}
