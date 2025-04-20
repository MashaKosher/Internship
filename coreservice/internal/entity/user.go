package entity

type User struct {
	Balance float64 `json:"balance"`
	Login   string  `json:"login"`
	WinRate float64 `json:"win-rate"`
}

type Response struct {
	Message string `json:"message"`
}

type BalanceBody struct {
	Balance float64 `json:"balance"`
}

type TypeResponse struct {
	User
	Message string `json:"message"`
}
