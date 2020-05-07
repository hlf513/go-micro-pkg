package kafka

import (
	"encoding/json"
	"sync"

	"github.com/hlf513/go-micro-pkg/config"
)

type Kafka struct {
	Host  []string `json:"host"`
	Topic string   `json:"topic"`
}

// kConf 定义更新配置
type kConf struct {
	Kafka *Kafka
}

var (
	kafkaConf = &kConf{}
	kl        sync.RWMutex
)

func GetConf() *Kafka {
	kl.RLock()
	defer kl.RUnlock()

	return kafkaConf.Kafka
}

// SetConf 更新配置
func SetConf(c []byte) error {
	kl.Lock()
	defer kl.Unlock()

	if c == nil {
		if err := config.GetConfigurator().Get([]string{"kafka"}, &kafkaConf.Kafka); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(c, kafkaConf); err != nil {
			return err
		}
	}

	return nil
}
