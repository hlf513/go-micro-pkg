package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/micro/go-micro/util/log"
)

// Setup 初始化 sentry
func Setup() {
	sentryConf, err := GetSentryConf()
	if err != nil {
		log.Fatal("[sentry setup] ", err.Error())
		return
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: sentryConf.Dns,
	}); err != nil {
		log.Fatal("[sentry setup] initialise failed:", err.Error())
	}
}
