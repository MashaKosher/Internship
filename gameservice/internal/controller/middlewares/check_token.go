package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gameservice/internal/adapter/kafka/consumers"
	"gameservice/internal/adapter/kafka/producers"
	"gameservice/internal/di"
)

func CheckTokenMiddleWare(logger di.LoggerType, cfg di.ConfigType, bus di.Bus) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("Middleware for Token check")

			// Получаем cookies
			accessToken, err := c.Cookie("access")
			if err != nil || accessToken.Value == "" {
				logger.Error("No access token")
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "No access token"})
			}
			logger.Info("Access Token:" + accessToken.Value)

			refreshToken, _ := c.Cookie("refresh")

			// Отправка токена в authserver через Kafka
			producers.CheckToken(accessToken.Value, refreshToken.Value, cfg, bus)

			// Получение ответа от authserver
			answer, _ := consumers.RecieveTokenInfo(cfg, bus)

			// Проверка валидности токена
			if answer.Err != "" {
				logger.Error(answer.Err)
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": answer.Err})
			}

			// Обновление access токена если нужно
			if answer.NewAccessToken != "" {
				logger.Info("New access Token")
				c.SetCookie(&http.Cookie{
					Name:     "access",
					Value:    answer.NewAccessToken,
					MaxAge:   3600,
					Path:     "/",
					HttpOnly: true,
					Secure:   false,
				})
			}

			c.Set("message", "Player added succesfully")
			c.Set("data", answer)
			return next(c)
		}
	}
}
