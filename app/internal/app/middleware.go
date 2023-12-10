package app

import (
	"github.com/go-chi/chi/v5"
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

func NewMetricsMiddleware(metrics *metrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
	}
}

func (m *MetricsMiddleware) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &StatusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(sw, r)
		pattern := chi.RouteContext(r.Context()).RoutePattern()
		m.metrics.latency.WithLabelValues(pattern, r.Method, strconv.Itoa(sw.statusCode)).Observe(time.Since(start).Seconds())
		m.metrics.rps.WithLabelValues(pattern, r.Method, strconv.Itoa(sw.statusCode)).Inc()
	}

	return http.HandlerFunc(fn)
}

type StatusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *StatusWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
