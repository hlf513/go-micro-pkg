package hystrix

import (
	"net"
	"net/http"

	hg "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/util/log"
)

func StreamServer() {
	conf := GetConf()
	if conf != nil && conf.StreamServer {
		hystrixStreamHandler := hg.NewStreamHandler()
		hystrixStreamHandler.Start()
		go func() {
			defer func() {
				if p := recover(); p != nil {
					log.Fatal("hystrix streaming server was failed:", p)
				}
			}()
			if err := http.ListenAndServe(net.JoinHostPort("", conf.ServerPort), hystrixStreamHandler); err != nil {
				panic(err)
			}
		}()
	}
}
