package middleware

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"fmt"
	"net/http"
)

func CheckToken(cfg di.ConfigType, bus di.Bus, logger di.LoggerType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Middlewware for token check")

			// Extracting Access Token
			accessToken, err := r.Cookie("access")
			if err != nil || (len(accessToken.String()) == 0) {
				http.Error(w, "No access token", http.StatusBadRequest)
				return
			}

			// Extracting Refresh Token
			refreshToken, err := r.Cookie("refresh")
			if err != nil || (len(accessToken.String()) == 0) {
				http.Error(w, "No refresh token", http.StatusBadRequest)
				return
			}

			// Executing Auth Response to AuthService
			bus.AuthProducer.CheckToken(accessToken.Value, refreshToken.Value)
			logger.Info("Auth Request sent")

			// Recieving Auth Response from AuthService
			var val entity.AuthAnswer
			val, _ = bus.AuthConsumer.AnswerTokens()
			logger.Info("Auth answer received: " + fmt.Sprint(val))

			if val.Role != "admin" {
				http.Error(w, "User is not admin", http.StatusForbidden)
				return
			}

			// If Refresh token is not valid field Err will be not empty
			if len(val.Err) != 0 {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}

			// If refresh token is valid, we recieve new access and set it in Cookie
			if len(val.NewAccessToken) != 0 {
				logger.Info("New access Token")
				http.SetCookie(w, &http.Cookie{
					Name:     "access",
					Value:    val.NewAccessToken,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   3600,
				})
			}

			next.ServeHTTP(w, r)
		})
	}
}
