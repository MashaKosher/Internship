package entity

import (
	"time"
)

type GameSettings struct {
	ID          uint          `json:"-" gorm:"primaryKey"`
	WinAmount   float64       `json:"win-amount"`
	LoseAmount  float64       `json:"lose-amount"`
	WaitingTime time.Duration `json:"waiting-time"`
}

func (settings *GameSettings) ToJSON() SettingsJson {
	return SettingsJson{
		LoseAmount:  settings.LoseAmount,
		WinAmount:   settings.WinAmount,
		WaitingTime: int(settings.WaitingTime.Seconds()),
	}
}

type SettingsJson struct {
	WinAmount   float64 `json:"win-amount" example:"10.05" validate:"min=0"`
	LoseAmount  float64 `json:"lose-amount" example:"5.97" validate:"min=0"`
	WaitingTime int     `json:"waiting-time" example:"3" validate:"min=0"`
}

func (settings *SettingsJson) ToDB() GameSettings {
	return GameSettings{
		LoseAmount:  settings.LoseAmount,
		WinAmount:   settings.WinAmount,
		WaitingTime: time.Second * time.Duration(settings.WaitingTime),
	}
}
