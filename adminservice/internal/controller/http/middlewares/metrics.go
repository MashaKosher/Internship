package middleware

import (
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/cpu"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_admin",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds_admin",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"path"},
	)

	memUsageMB = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_mb_admin",
			Help: "Current memory usage in MB",
		},
	)
	cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent_admin",
			Help: "CPU usage percentage",
		},
	)
)

func UpdateMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memUsageMB.Set(float64(m.Alloc) / (1024 * 1024))

	percent, _ := cpu.Percent(0, false)
	cpuUsage.Set(percent[0])
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		// Обертка для получения статуса ответа
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(path).Observe(duration)
		httpRequestsTotal.WithLabelValues(
			r.Method,
			path,
			http.StatusText(rw.status),
		).Inc()

	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
