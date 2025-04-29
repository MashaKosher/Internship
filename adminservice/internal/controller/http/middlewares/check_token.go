package middleware

import (
	"adminservice/internal/adapter/kafka/consumers"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"context"
	"log"
	"net/http"
)

func CheckToken(cfg di.ConfigType, bus di.Bus) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Middlewware for token check")

			// Extracting Access Token
			accessToken, err := r.Cookie("access")
			if err != nil || (len(accessToken.String()) == 0) {
				http.Error(w, "No token", http.StatusBadRequest)
				return
			}

			log.Println(accessToken.Value)

			// Extracting Refresh Token
			refreshToken, _ := r.Cookie("refresh")

			// Executing Auth Response to AuthService
			producers.CheckToken(accessToken.Value, refreshToken.Value, cfg, bus)
			log.Println("Auth Request sent")

			// Recieving Auth Response from AuthService
			var val entity.AuthAnswer
			val, _ = consumers.AnswerTokens(cfg, bus)
			log.Println("Auth answer received")
			log.Println(val)

			if val.Role != "admin" {
				http.Error(w, "User is not admin", http.StatusBadRequest)
				return
			}

			// If Refresh token is not valid field Err will be not empty
			if len(val.Err) != 0 {
				http.Error(w, "Access Token expired", http.StatusBadRequest)
				return
			}

			// If refresh token is valid, we recieve new access and set it in Cookie
			if len(val.NewAccessToken) != 0 {
				log.Println("New access Token")
				http.SetCookie(w, &http.Cookie{
					Name:     "access",
					Value:    val.NewAccessToken,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   3600,
				})
			}

			ctx := context.WithValue(r.Context(), "val", val)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Middlewware for token check")

// 		// Extracting Access Token
// 		accessToken, err := r.Cookie("access")
// 		if err != nil || (len(accessToken.String()) == 0) {
// 			http.Error(w, "No token", http.StatusBadRequest)
// 			return
// 		}

// 		log.Println(accessToken.Value)

// 		// Extracting Refresh Token
// 		refreshToken, _ := r.Cookie("refresh")

// 		// Executing Auth Response to AuthService
// 		producers.CheckToken(accessToken.Value, refreshToken.Value, cfg, bus)
// 		log.Println("Auth Request sent")

// 		// Recieving Auth Response from AuthService
// 		var val entity.AuthAnswer
// 		val, _ = consumers.AnswerTokens()
// 		log.Println("Auth answer received")
// 		log.Println(val)

// 		if val.Role != "admin" {
// 			http.Error(w, "User is not admin", http.StatusBadRequest)
// 			return
// 		}

// 		// If Refresh token is not valid field Err will be not empty
// 		if len(val.Err) != 0 {
// 			http.Error(w, "Access Token expired", http.StatusBadRequest)
// 			return
// 		}

// 		// If refresh token is valid, we recieve new access and set it in Cookie
// 		if len(val.NewAccessToken) != 0 {
// 			log.Println("New access Token")
// 			http.SetCookie(w, &http.Cookie{
// 				Name:     "access",
// 				Value:    val.NewAccessToken,
// 				Path:     "/",
// 				HttpOnly: true,
// 				MaxAge:   3600,
// 			})
// 		}

// 		ctx := context.WithValue(r.Context(), "val", val)
// 		r = r.WithContext(ctx)

// 		next.ServeHTTP(w, r)
// 	})
// }
