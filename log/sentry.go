package log

import (
	"time"

	"github.com/getsentry/sentry-go"
)

const sentryFlushTimeout = 3 * time.Second

// SentryMessage 只会记录 msg，不会记录 stack
func SentryMessage(msg string) {
	sentry.CaptureMessage(msg)
	sentry.Flush(sentryFlushTimeout)
}

// SentryMessage 记录 msg 和 stack
func SentryException(err error) {
	sentry.CaptureException(err)
	sentry.Flush(sentryFlushTimeout)
}
