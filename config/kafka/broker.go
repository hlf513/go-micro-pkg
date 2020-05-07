package kafka

import (
	"fmt"

	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/kafka/v2"
)

var kafkaBroker broker.Broker

func Broker() (broker.Broker, error) {
	if kafkaBroker == nil {
		kafkaBroker = kafka.NewBroker(func(options *broker.Options) {
			options.Addrs = GetConf().Host
		})

		if err := kafkaBroker.Init(); err != nil {
			return nil, fmt.Errorf("kafka broker init exception：%s", err.Error())
		}

		if err := kafkaBroker.Connect(); err != nil {
			return nil, fmt.Errorf("kafka broker connect exception：%s", err.Error())
		}
	}

	return kafkaBroker, nil
}

func Close() {
	if kafkaBroker != nil {
		_ = kafkaBroker.Disconnect()
	}
}
