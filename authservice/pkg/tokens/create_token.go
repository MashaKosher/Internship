package tokens

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(tokenType di.TokenType, user *entity.User, logger di.LoggerType, keys di.RSAKeys) (string, error) {
	logger.Info("Creating " + string(tokenType) + " token " + " for User with ID " + fmt.Sprint(user.ID))
	privateKey := keys.PrivateKey
	method := jwt.SigningMethodRS256
	logger.Info("Signing Method: " + method.Name + "with Private key ")

	// Filling payload
	claims := fillClaims(tokenType, user)

	// Generating token
	token, err := jwt.NewWithClaims(method, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func fillClaims(tokenType di.TokenType, user *entity.User) jwt.MapClaims {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
	}
	if tokenType == di.ACCESS_TOKEN {
		claims["username"] = user.Username
		claims["type"] = string(di.ACCESS_TOKEN)
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		// claims["exp"] = time.Now().Add(time.Second * 30).Unix()
		// claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	} else {
		claims["type"] = string(di.REFRESH_TOKEN)
		// claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
		// claims["exp"] = time.Now().Add(time.Minute).Unix()
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	}
	return claims
}
