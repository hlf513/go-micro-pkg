package ginTrace

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/hlf513/go-micro-pkg/config/jaeger"
)

// RecoveryWrapper 恢复 panic 并记录到日志中
func RecoveryWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				span := opentracing.SpanFromContext(c.Request.Context())
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				errorMsg := fmt.Sprintf("%v", err)

				if brokenPipe {
					if span != nil {
						jaeger.SetError(span, errors.New(errorMsg))
					}
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				errorMsg = fmt.Sprintf(`[Recovery from panic] - %s - %s`, errorMsg, debug.Stack())
				if span != nil {
					jaeger.SetError(span, errors.New(errorMsg))
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
