package entity

import "time"

type GameSettings struct {
	ID          uint          `json:"-" gorm:"primaryKey"`
	WinAmount   float64       `json:"win-amount"`
	LoseAmount  float64       `json:"lose-amount"`
	WaitingTime time.Duration `json:"waiting-time"`
}

type SettingsJson struct {
	WinAmount   float64 `json:"win-amount"`
	LoseAmount  float64 `json:"lose-amount"`
	WaitingTime int     `json:"waiting-time"`
}
