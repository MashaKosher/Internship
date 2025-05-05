package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func IntiMetricsRoutes(r *chi.Mux) {
	r.Handle("/metrics", promhttp.Handler())
}
