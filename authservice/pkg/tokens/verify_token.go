package tokens

import (
	"authservice/internal/di"
	"errors"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt"
)

func TokenVerify(token string, logger di.LoggerType, keys di.RSAKeys) (*jwt.Token, error) {
	publicKey := keys.PublicKey
	validatedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.GetSigningMethod("RS256").Alg() {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})

	logger.Info("\nToken is valid need to check time\n")

	if exp, ok := validatedToken.Claims.(jwt.MapClaims)["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		logger.Info(fmt.Sprint("Now: ", time.Now(), " Exp: ", expirationTime))
		if time.Now().After(expirationTime) {
			return nil, errors.New("token expired")
		}

	}

	return validatedToken, err
}

func VerifyTokenType(must, actually string) error {
	if must != actually {
		return errors.New("Token must have type " + must)
	}
	return nil
}
