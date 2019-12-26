package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
)

// 默认上报超时时间 1s
const sentryFlushTimeout = 1 * time.Second

// Message 只会记录 msg，不会记录 stack
func Message(msg string) {
	sentry.CaptureMessage(msg)
	sentry.Flush(sentryFlushTimeout)
}

// Exception 记录 msg 和 stack
func Exception(err error) {
	sentry.CaptureException(err)
	sentry.Flush(sentryFlushTimeout)
}
