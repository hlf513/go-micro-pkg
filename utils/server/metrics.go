package server

import (
	"net/http"

	"github.com/micro/go-micro/v2/util/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus 启动 Metrics 服务
func Prometheus(addr string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal("Metrics Server:", err)
		}
	}()
	log.Info("Metrics [http] Listening on [::]" + addr)
}
