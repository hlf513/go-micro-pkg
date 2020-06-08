package async

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/hlf513/go-micro-pkg/config/jaeger"
	"github.com/hlf513/go-micro-pkg/config/zap"
	"github.com/hlf513/go-micro-pkg/wrapper/common"
)

// Handler 异步处理函数
// f 函数定义
// traceOperationName 请求追踪名
// funcParams f函数的入参(可不填)
func Handler(f func(ctx context.Context, params ...interface{}) error, traceOperationName string, funcParams ...interface{}) {

	// 平滑关闭
	wg := common.WaitGroup()
	wg.Add(1)
	defer wg.Done()

	// 记录tracing
	span, c := jaeger.StartTrace(traceOperationName)
	defer jaeger.StopTrace(span)
	defer func() {
		if p := recover(); p != nil {
			zap.Error(c,
				fmt.Sprintf(
					`[Recovery from panic] - %s`,
					debug.Stack(),
				),
			)
		}
	}()

	if err := f(c, funcParams...); err != nil {
		jaeger.SetError(span, err)
		zap.Error(c, err.Error())
	}
}
