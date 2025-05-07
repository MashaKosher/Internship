package entity

import db "coreservice/internal/repository/sqlc/generated"

type User struct {
	Balance float64 `json:"balance"`
	Login   string  `json:"login"`
	WinRate float64 `json:"win-rate"`
}

type SignedUpUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	ReferalID int    `json:"referal-id"`
}

func (user *SignedUpUser) ToAuthAnswer() AuthAnswer {
	return AuthAnswer{
		ID:    int32(user.ID),
		Login: user.Username,
	}
}

type Balance struct {
	Balance float64 `json:"balance" example:"10.20" validate:"gte=0"`
}

type Response struct {
	Message string `json:"message"`
}

type TypeResponse struct {
	User db.User

	Message string `json:"message"`
}
