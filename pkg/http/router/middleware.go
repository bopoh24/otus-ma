package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"time"
)

// LatencyRpsMetrics is a middleware that collects metrics about requests:
// - latency
// - rps

type MetricsMiddleware struct {
	metrics *metrics
}

// NewMetricsMiddleware returns a new MetricsMiddleware.
func NewMetricsMiddleware(namespace string) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: newMetrics(namespace),
	}
}

func (m *MetricsMiddleware) Metrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(sw, r)
		pattern := chi.RouteContext(r.Context()).RoutePattern()
		// skip healthz, readyz and metrics endpoints
		if pattern == "/healthz" || pattern == "/readyz" || pattern == "/metrics" {
			return
		}
		m.metrics.latency.WithLabelValues(pattern, r.Method, strconv.Itoa(sw.statusCode)).Observe(time.Since(start).Seconds())
		m.metrics.rps.WithLabelValues(pattern, r.Method, strconv.Itoa(sw.statusCode)).Inc()
	}

	return http.HandlerFunc(fn)
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type metrics struct {
	latency *prometheus.HistogramVec
	rps     *prometheus.CounterVec
}

func newMetrics(namespace string) *metrics {
	return &metrics{
		latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "http request latency",
			Buckets:   []float64{0.1, 0.2, 0.5, 1, 2, 5},
		}, []string{"url", "method", "status"}),
		rps: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "http request count",
		}, []string{"url", "method", "status"}),
	}
}
