package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
)

func TracingMiddleware(c *fiber.Ctx) error {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan(c.Path())
	defer span.Finish()

	span.SetTag("http.method", c.Method())
	span.SetTag("http.url", c.Path())

	ctx := opentracing.ContextWithSpan(c.Context(), span)
	c.Locals("tracing_ctx", ctx)

	return c.Next()
}
