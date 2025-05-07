package middleware

import (
	"net/http"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/cpu"
)

// Метрики
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_auth",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds_auth",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)

	memUsageMB = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_mb_auth",
			Help: "Current memory usage in MB",
		},
	)
	cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent_auth",
			Help: "CPU usage percentage",
		},
	)
)

func UpdateMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memUsageMB.Set(float64(m.Alloc) / (1024 * 1024)) // Преобразуем байты в МБ

	percent, _ := cpu.Percent(0, false)
	cpuUsage.Set(percent[0])
}

// Middleware для сбора метрик
func MetricsMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/metrics" {
		return c.Next()
	}

	method := c.Method()
	path := c.Path()

	// duration := time.Since(start).Seconds()
	// httpRequestDuration.WithLabelValues(path).Observe(duration)
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(method, path))
	defer timer.ObserveDuration()

	httpRequestsTotal.WithLabelValues(
		method,
		path,
		http.StatusText(c.Response().StatusCode()),
	).Inc()

	return c.Next()
}
