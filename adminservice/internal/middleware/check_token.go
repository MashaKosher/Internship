package middleware

import (
	"adminservice/internal/adapter/kafka/consumers"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middlewware for token check")

		accessToken, err := r.Cookie("access")
		if err != nil || (len(accessToken.String()) == 0) {
			http.Error(w, "No token", http.StatusBadRequest)
			return
		}

		log.Println(accessToken.Value)

		refreshToken, _ := r.Cookie("refresh")

		var val entity.AuthAnswer
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			producers.CheckToken(accessToken.Value, refreshToken.Value)
			fmt.Println("Request sent")
		}()

		wg.Wait()

		val, _ = consumers.AnswerTokens()
		fmt.Println("Answer received")
		log.Println(val)

		// If token is not valid fiels Err will be not empty
		if len(val.Err) != 0 {
			http.Error(w, "Access Token expired", http.StatusBadRequest)
			return
		}

		// If refresh token is valid, we recieve new access
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
