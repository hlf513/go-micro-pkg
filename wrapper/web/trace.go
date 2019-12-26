package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/util/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// GinTraceWrapper 注入 metadata 到 tracing（原生 micro srv 的 trace 依赖 metadata）
func GinTraceWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		md := make(map[string]string)
		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
		if err != nil {
			log.Fatal("[jaeger metadata]:", err.Error())
		}
		// ctx := context.TODO()
		ctx := c.Request.Context()
		ctx = metadata.NewContext(ctx, md)
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// bodyLogWriter 暂存响应
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 输出响应
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinTraceLogWrapper 记录 gin 的请求、响应、错误日志
func GinTraceLogWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		if span == nil {
			c.Next()
			return
		}

		// start := time.Now()
		// 缓存 response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// request log
		type req struct {
			ContextType string `json:"context_type"`
			Method      string
			Url         string
			Params      string
		}
		var r req
		r.Method = c.Request.Method
		r.Url = c.Request.URL.String()
		r.ContextType = c.ContentType()

		if c.ContentType() == "application/json" { // 记录 json 请求
			body, _ := c.GetRawData()
			// here is important
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			r.Params = string(body)
		} else { // 记录其他请求
			var p = make(map[string]interface{})
			_ = c.Request.ParseMultipartForm(1 << 30)
			for k, v := range c.Request.PostForm {
				p[k] = v
			}
			param, _ := json.Marshal(p)
			r.Params = string(param)
		}

		msg, _ := json.Marshal(r)
		span.LogKV("request", string(msg))

		c.Next()

		// end := time.Now()
		// latency := end.Sub(start)

		// error log
		if len(c.Errors) > 0 {
			for n, e := range c.Errors.Errors() {
				ext.SamplingPriority.Set(span, 1)
				span.SetTag("error", true)
				span.LogKV("error_msg_"+strconv.Itoa(n), e)
			}
		}

		// response log
		span.LogKV("response", blw.body.String())
	}
}

// GinRecoveryWithZap 恢复 panic 并记录到日志中
func GinRecoveryWithZap() gin.HandlerFunc {
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
						ext.SamplingPriority.Set(span, 1)
						span.SetTag("error", true)
						span.LogKV("error_msg", errorMsg)
					}
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				errorMsg = fmt.Sprintf(`[Recovery from panic] - %s - %s`, errorMsg, debug.Stack())
				if span != nil {
					ext.SamplingPriority.Set(span, 1)
					span.SetTag("error", true)
					span.LogKV("error_msg", errorMsg)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
