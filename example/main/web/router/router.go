package router

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hlf513/go-micro-pkg/config/server"
	"github.com/hlf513/go-micro-pkg/wrapper/common"
)

func Register(app *gin.Engine) {
	// 拦截非法请求
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, "404")
	})

	// 定义 group prefix；和 server name 保持一致
	serverNames := strings.Split(server.GetConf().Name, ".")
	apiPrefix := "/" + serverNames[len(serverNames)-1]
	g := app.Group(apiPrefix)

	// example
	g.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, "ok")
		// c := ctx.Copy()
		go func() {
			// 只有使用 common.WaitGroup() 才能平滑关闭
			common.WaitGroup().Add(1)	
			defer common.WaitGroup().Done()
			
			
		}()
		
	})
}
