package entity

type UserInDTO struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type LoginOutDTO struct {
	UserID       int    `json:"id"`
	UserName     string `json:"username"`
	UserRole     string `json:"role"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}
