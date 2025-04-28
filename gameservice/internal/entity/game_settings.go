package entity

type GameSettings struct {
	WinAmount   float64 `json:"win-amount"`
	LoseAmount  float64 `json:"lose-amount"`
	WaitingTime int     `json:"waiting-time"`
}
