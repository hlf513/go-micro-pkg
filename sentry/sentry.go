package sentry

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/micro/go-micro/util/log"

	"github.com/hlf513/go-micro-pkg/config"
)

func Init(s config.Sentry) {
	config.SetSentry(s)
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: s.Dns,
	}); err != nil {
		log.Fatal(fmt.Sprintf("Sentry initialization failed: %v\n", err), "application")
	}
}
