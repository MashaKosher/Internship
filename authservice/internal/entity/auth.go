package entity

type AuthAnswer struct {
	Role           string `json:"role"`
	ID             int32  `json:"id"`
	Login          string `json:"login"`
	Err            string `json:"err"`
	NewAccessToken string `json:"new-access-token"`
}

type AuthRequest struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
	Partition    int32  `json:"partition"`
}

type Password struct {
	NewPassword string `json:"new-password"`
}

type Response struct {
	Message string `json:"message"`
}
