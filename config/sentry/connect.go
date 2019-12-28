package sentry

import (
	client "github.com/getsentry/sentry-go"
	"github.com/micro/go-micro/util/log"
)

// Connect 连接 sentry server
func Connect() {
	sentryConf := GetConf()

	if err := client.Init(client.ClientOptions{
		Dsn: sentryConf.Dns,
	}); err != nil {
		log.Fatal("[sentry setup] initialise failed:", err.Error())
	}
}
