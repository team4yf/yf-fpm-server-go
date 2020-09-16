package fpm

import "github.com/prometheus/client_golang/prometheus"

var (
	bizExecuteVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "biz",
			Subsystem: "modules",
			Name:      "execute_total",
			Help:      "Total number of biz executed",
		},
		[]string{"method", "result"},
	)
)

func registerPrometheus(fpmApp *Fpm) {
	prometheus.MustRegister(bizExecuteVec)
	fpmApp.routers.Handle("/metrics", prometheus.Handler())
}

func incBizExecuteVec(method, result string) {
	bizExecuteVec.WithLabelValues(method, result).Inc()
}
