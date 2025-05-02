package entity

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

type UserInDTO struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type UserOutDTO struct {
	UserID       int    `json:"id"`
	UserName     string `json:"username"`
	UserRole     string `json:"role"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (u *User) ToDTO(accessToken, refreshToken string) UserOutDTO {
	return UserOutDTO{UserID: int(u.ID), UserName: u.Username, UserRole: string(u.Role), AccessToken: accessToken, RefreshToken: refreshToken}
}
