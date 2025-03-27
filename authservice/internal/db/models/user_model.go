package models

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID       uint   `json:"-" gorm:"primaryKey"`
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
	Role     Role   `json:"-" gorm:"default:user"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Role         Role   `json:"role"`
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
	TokenType    string `json:"token-type"`
}
