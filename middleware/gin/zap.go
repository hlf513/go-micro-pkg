package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hlf513/go-micro-pkg/log"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Ginzap 记录请求和响应日志
func Ginzap() gin.HandlerFunc {
	return func(c *gin.Context) {
		// start := time.Now()
		// cache response body
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

		if c.ContentType() == "application/json" {
			body, _ := c.GetRawData()
			// important
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			r.Params = string(body)
		} else {
			var p = make(map[string]interface{})
			_ = c.Request.ParseMultipartForm(1 << 30)
			for k, v := range c.Request.PostForm {
				p[k] = v
			}
			param, _ := json.Marshal(p)
			r.Params = string(param)
		}

		msg, _ := json.Marshal(r)
		log.Info(c.Request.Context(), string(msg), "request")

		c.Next()

		// end := time.Now()
		// latency := end.Sub(start)

		// error log
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Error(c.Request.Context(), e, "response")
			}
		}

		// response log
		log.Info(c.Request.Context(), blw.body.String(), "response")
	}
}

// RecoveryWithZap 记录 recovery 日志
func RecoveryWithZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
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

				if brokenPipe {
					errorMsg := fmt.Sprintf("%v", err)
					log.Error(c.Request.Context(), errorMsg, "request")
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				log.Error(c.Request.Context(), "[Recovery from panic]", "request")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
