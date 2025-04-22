package auth

type UserInDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutDTO struct {
	UserID       int    `json:"id"`
	UserName     string `json:"username"`
	UserRole     string `json:"role"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}
