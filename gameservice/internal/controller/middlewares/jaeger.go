package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

func TracingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan(c.Path())
		defer span.Finish()

		span.SetTag("http.method", c.Request().Method)
		span.SetTag("http.url", c.Request().URL.Path)

		ctx := opentracing.ContextWithSpan(c.Request().Context(), span)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		// Якщо потрібно передати контекст у handler:
		c.Set("tracing_ctx", ctx)

		return next(c)
	}
}
