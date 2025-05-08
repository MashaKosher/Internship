package tokens

import (
	"authservice/internal/entity"

	"github.com/golang-jwt/jwt"
)

func GetIdFromValidatedToken(validatedToken *jwt.Token) (int, error) {
	unvalidatedUserID := validatedToken.Claims.(jwt.MapClaims)["sub"]

	validatedUserID, ok := unvalidatedUserID.(float64)
	if !ok {
		return -1, entity.ErrInavlidID
	}

	return int(validatedUserID), nil
}

func GetTypeFromValidatedToken(validatedToken *jwt.Token) (string, error) {
	unvalidatedTokenType := validatedToken.Claims.(jwt.MapClaims)["type"]

	validatedTokenType, ok := unvalidatedTokenType.(string)
	if !ok {
		return "", entity.ErrInvalidTokenType
	}

	return validatedTokenType, nil
}
