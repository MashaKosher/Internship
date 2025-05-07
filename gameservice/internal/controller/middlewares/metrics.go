package middlewares

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/cpu"
)

// Создаем и регистрируем метрики
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_game",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds_game",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)

	memUsageMB = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_mb_game",
			Help: "Current memory usage in MB",
		},
	)
	cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent_game",
			Help: "CPU usage percentage",
		},
	)
)

func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		path := c.Path()

		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(method, path))
		defer timer.ObserveDuration()

		err := next(c)

		status := http.StatusText(c.Response().Status)
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()

		return err
	}
}

func UpdateMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memUsageMB.Set(float64(m.Alloc) / (1024 * 1024)) // Преобразуем байты в МБ

	percent, _ := cpu.Percent(0, false)
	cpuUsage.Set(percent[0])
}
