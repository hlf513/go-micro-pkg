package app

import (
	"github.com/gin-gonic/gin"

	"go-micro-pkg/jaeger"
	"go-micro-pkg/log"
)

type Response struct {
	C *gin.Context
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Response) Success(data interface{}, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = GetMsg(Success)
	}
	g.C.JSON(200, response{
		Code: Success,
		Msg:  message,
		Data: data,
	})
}

func (g *Response) Error(httpCode int, errCode int, message string, data interface{}) {
	log.Trace(g.C.Request.Context()).Sampler()
	g.C.JSON(httpCode, response{
		Code: errCode,
		Msg:  message,
		Data: data,
	})
}

func (g *Response) Fail(errCode int, errMsg string, message ...string) {
	if errMsg != "" {
		log.Trace(g.C.Request.Context()).Sampler()
		log.Error(g.C.Request.Context(), errMsg, "response")
	}
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = GetMsg(errCode)
	}
	g.C.JSON(200, response{
		Code: errCode,
		Msg:  msg,
		Data: jaeger.GetTraceId(g.C.Request.Context()),
	})
}
