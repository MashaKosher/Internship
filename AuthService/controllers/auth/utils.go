package auth

import (
	config "authservice/config"
	models "authservice/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passwrod string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwrod), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func ValidatePassword(HashedPass, RawPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(RawPass))
}

func GenerateToken(user *models.User) (string, error) {

	privateKey := config.RSAkeys.PrivateKey
	method := jwt.SigningMethodRS256
	config.Logger.Info("Signing Method: " + method.Name + " Private key: " + fmt.Sprint(privateKey))

	// Filling payload
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 1).Unix(),
	}

	// Generating token
	token, err := jwt.NewWithClaims(method, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func TokenValidation(token string) (*jwt.Token, error) {
	publicKey := config.RSAkeys.PublicKey
	validatedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.GetSigningMethod("RS256").Alg() {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})

	if exp, ok := validatedToken.Claims.(jwt.MapClaims)["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return nil, errors.New("token expired")
		}
	}

	return validatedToken, err
}

func ConvertUserToResponse(token string, user *models.User) models.UserResponse {
	return models.UserResponse{
		Role:      user.Role,
		ID:        user.ID,
		Username:  user.Username,
		Token:     token,
		TokenType: "Bearer",
	}
}
