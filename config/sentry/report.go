package sentry

import (
	"time"

	client "github.com/getsentry/sentry-go"
)

// 默认上报超时时间 1s
const sentryFlushTimeout = 1 * time.Second

// Message 只会记录 msg，不会记录 stack
func Message(msg string) {
	client.CaptureMessage(msg)
	client.Flush(sentryFlushTimeout)
}

// Exception 记录 msg 和 stack
func Exception(err error) {
	client.CaptureException(err)
	client.Flush(sentryFlushTimeout)
}
