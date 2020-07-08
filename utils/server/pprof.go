package server

import (
	"net/http"
	// 开启 pprof
	_ "net/http/pprof"

	"github.com/micro/go-micro/v2/util/log"
)

// Server pprof http server
func Server(addr string) {
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal("Profile Server:", err)
		}
	}()
	log.Info("pprof [http] Listening on [::]" + addr)
}
