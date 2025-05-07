package middleware

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_core",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_core",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)

	memUsageMB = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_mb_core",
			Help: "Current memory usage in MB",
		},
	)
	cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent_core",
			Help: "CPU usage percentage",
		},
	)
)

func UpdateMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memUsageMB.Set(float64(m.Alloc) / (1024 * 1024)) // Конвертация в MB

	// Использование CPU можно получать через gopsutil
	percent, _ := cpu.Percent(0, false)
	cpuUsage.Set(percent[0])
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Получаем путь (с fallback)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Замер времени и инкремент счетчика
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		))
		defer timer.ObserveDuration()

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
		).Inc()

		c.Next()
	}
}
