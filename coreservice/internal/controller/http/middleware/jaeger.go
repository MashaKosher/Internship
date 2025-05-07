package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

// TracingMiddleware створює span для кожного HTTP-запиту в Gin
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan(c.FullPath()) // або c.Request.URL.Path
		defer span.Finish()

		span.SetTag("http.method", c.Request.Method)
		span.SetTag("http.url", c.Request.URL.Path)

		ctx := opentracing.ContextWithSpan(c.Request.Context(), span)
		c.Request = c.Request.WithContext(ctx)

		// Якщо потрібно дістати в handler'ах:
		c.Set("tracing_ctx", ctx)

		c.Next()
	}
}
