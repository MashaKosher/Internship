package middleware

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan(r.URL.Path)
		defer span.Finish()

		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.URL.String())

		ctx := opentracing.ContextWithSpan(r.Context(), span)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
