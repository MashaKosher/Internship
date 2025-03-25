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

// type TokenSettings struct {
// 	ACCESS_TOKEN  string
// 	REFRESH_TOKEN string
// }

const ACCESS_TOKEN = "access"
const REFRESH_TOKEN = "refresh"

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

func TokenVerify(token string) (*jwt.Token, error) {
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

func ConvertUserToResponse(accessToken string, refreshToken string, user *models.User) models.UserResponse {
	return models.UserResponse{
		Role:         user.Role,
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}
}

func GetIdFromValidatedToken(validatedToken *jwt.Token) (int, error) {
	unvalidatedUserID := validatedToken.Claims.(jwt.MapClaims)["sub"]

	validatedUserID, ok := unvalidatedUserID.(float64)
	if !ok {
		return -1, errors.New("invalid ID value:" + fmt.Sprint(unvalidatedUserID))
	}

	return int(validatedUserID), nil
}

func GetTypeFromValidatedToken(validatedToken *jwt.Token) (string, error) {
	unvalidatedTokenType := validatedToken.Claims.(jwt.MapClaims)["type"]

	validatedTokenType, ok := unvalidatedTokenType.(string)
	if !ok {
		return "", errors.New("invalid Token type:" + fmt.Sprint(unvalidatedTokenType))
	}

	return validatedTokenType, nil
}

func CreateToken(tokenType string, user *models.User) (string, error) {
	config.Logger.Info("Creating " + tokenType + " token " + " for User with ID " + fmt.Sprint(user.ID))
	privateKey := config.RSAkeys.PrivateKey
	method := jwt.SigningMethodRS256
	config.Logger.Info("Signing Method: " + method.Name + "with Private key ")

	// Filling payload
	claims := fillClaims(tokenType, user)

	// Generating token
	token, err := jwt.NewWithClaims(method, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func fillClaims(tokenType string, user *models.User) jwt.MapClaims {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
	}
	if tokenType == ACCESS_TOKEN {
		claims["username"] = user.Username
		claims["type"] = ACCESS_TOKEN
		claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	} else {
		claims["type"] = REFRESH_TOKEN
		claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	}
	return claims
}

func VerifyTokenType(must, actually string) error {
	if must != actually {
		return errors.New("Token must have type " + must)
	}
	return nil
}
