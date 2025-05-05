package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Метрики
// var (
// 	httpRequestsTotal = promauto.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total number of HTTP requests",
// 		},
// 		[]string{"method", "path", "status"},
// 	)

// 	httpDuration = promauto.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_response_time_seconds",
// 			Help:    "Duration of HTTP requests",
// 			Buckets: []float64{0.1, 0.5, 1, 2, 5},
// 		},
// 		[]string{"path"},
// 	)
// )

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"path"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		// Обертка для получения статуса ответа
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues(path).Observe(duration)
		httpRequestsTotal.WithLabelValues(
			r.Method,
			path,
			http.StatusText(rw.status),
		).Inc()
	})
}
