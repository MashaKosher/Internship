package tokens

import (
	"authservice/internal/entity"
	"authservice/internal/keys"
	"authservice/pkg/logger"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const ACCESS_TOKEN = "access"
const REFRESH_TOKEN = "refresh"

func CreateToken(tokenType string, user *entity.User) (string, error) {
	logger.Logger.Info("Creating " + tokenType + " token " + " for User with ID " + fmt.Sprint(user.ID))
	privateKey := keys.RSAkeys.PrivateKey
	method := jwt.SigningMethodRS256
	logger.Logger.Info("Signing Method: " + method.Name + "with Private key ")

	// Filling payload
	claims := fillClaims(tokenType, user)

	// Generating token
	token, err := jwt.NewWithClaims(method, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func fillClaims(tokenType string, user *entity.User) jwt.MapClaims {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
	}
	if tokenType == ACCESS_TOKEN {
		claims["username"] = user.Username
		claims["type"] = ACCESS_TOKEN
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		// claims["exp"] = time.Now().Add(time.Second * 30).Unix()
		// claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	} else {
		claims["type"] = REFRESH_TOKEN
		// claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
		// claims["exp"] = time.Now().Add(time.Minute).Unix()
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	}
	return claims
}
