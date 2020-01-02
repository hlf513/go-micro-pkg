package ginTrace

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)



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

// RecordLog 记录 gin 的请求、响应、错误日志
func RecordLog() gin.HandlerFunc {
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
				span.LogKV("error_"+strconv.Itoa(n), e, "debug.stack", string(debug.Stack()))
			}
		}

		// response log
		span.LogKV("response", blw.body.String())
	}
}

