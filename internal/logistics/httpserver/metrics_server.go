package httpserver

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsServer(httpAddr string, reg *prometheus.Registry) *http.Server {
	httpSrv := &http.Server{Addr: httpAddr}
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	httpSrv.Handler = m
	return httpSrv
}
