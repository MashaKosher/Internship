package v1

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/valyala/fasthttp/fasthttpadaptor"
// )

// func InitMetricsRoutes(app *fiber.App) {
// 	app.Get("/metrics", func(c *fiber.Ctx) error {
// 		// адаптируем http.Handler к fasthttp
// 		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.Context())
// 		return nil
// 	})
// }
