package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	latency *prometheus.HistogramVec
	rps     *prometheus.CounterVec
}

func newMetrics() *metrics {
	return &metrics{
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "user",
			Name:      "http_request_duration_seconds",
			Help:      "http request latency",
			Buckets:   []float64{0.1, 0.2, 0.5, 1, 2, 5},
		}, []string{"url", "method", "status"}),
		rps: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "user",
			Name:      "http_requests_total",
			Help:      "http request count",
		}, []string{"url", "method", "status"}),
	}
}
