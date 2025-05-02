package entity

import "time"

type GameSettings struct {
	ID          uint          `json:"-" gorm:"primaryKey"`
	WinAmount   float64       `json:"win-amount"`
	LoseAmount  float64       `json:"lose-amount"`
	WaitingTime time.Duration `json:"waiting-time"`
}

type SettingsJson struct {
	WinAmount   float64 `json:"win-amount" validate:"min=0"`
	LoseAmount  float64 `json:"lose-amount" validate:"min=0"`
	WaitingTime int     `json:"waiting-time" validate:"min=0"`
}

func (settings *SettingsJson) ToDB() GameSettings {
	return GameSettings{
		LoseAmount:  settings.LoseAmount,
		WinAmount:   settings.WinAmount,
		WaitingTime: time.Duration(settings.WaitingTime),
	}
}
