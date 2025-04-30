package entity

type User struct {
	Balance float64 `json:"balance"`
	Login   string  `json:"login"`
	WinRate float64 `json:"win-rate"`
}

type Balance struct {
	Balance float64 `json:"balance" example:"10.20" validate:"gte=0"`
}

///////////////////////////////

type Response struct {
	Message string `json:"message"`
}

type TypeResponse struct {
	User
	Message string `json:"message"`
}
