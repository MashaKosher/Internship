package auth

import (
	models "authservice/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passwrod string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwrod), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func GenerateToken(user *models.User) (string, error) {
	// Becoming a secret
	godotenv.Load()
	secret := []byte(os.Getenv("JWT_SECRET"))
	method := jwt.SigningMethodHS256

	// method := jwt.SigningMethodRS256

	// Filling payload
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 2).Unix(),
	}

	// Generating token
	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ConvertUserToResponse(token string, user *models.User) models.UserResponse {
	return models.UserResponse{
		Role:     user.Role,
		ID:       user.ID,
		Username: user.Username,
		Token:    token,
	}
}

func ValidatePassword(HashedPass, RawPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(RawPass))
}

func TokenValidation(token string) (*jwt.Token, error) {
	godotenv.Load()
	secret := []byte(os.Getenv("JWT_SECRET"))
	restoken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
			return nil, errors.New("Invalid signing method")
		}
		return secret, nil
	})

	return restoken, err
}

// func CheckToken(token string) (int, error){
// 	godotenv.Load()
// 	secret := []byte(os.Getenv("JWT_SECRET"))
// 	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
// 		if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
// 			return nil, fiber.NewError(400, "Unexpceted signing method")
// 		}

// 		return secret, nil
// 	})

// 	// if token is invalid we clear the cookie
// 	if err != nil {
// 		c.ClearCookie()
// 		return fiber.NewError(400, "Invalid Token")
// 	}

// 	userId := token.Claims.(jwt.MapClaims)["userId"]

// 	if err := db.DB.Model(&models.User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
// 		c.ClearCookie()
// 		return fiber.NewError(400, "No such User")
// 	}

// }
