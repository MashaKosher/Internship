package entity

import "time"

type DBDailyTasks struct {
	TaskDate       time.Time `json:"-" gorm:"primaryKey;type:date" validate:"required"`
	ReferalsAmount int       `json:"referals-amount" example:"10" minimum:"0" validate:"gte=0"`
	GamesAmount    int       `json:"games-amount" example:"5" minimum:"0" validate:"gte=0"`
}

type DailyTasks struct {
	TaskDate       string `json:"task-date" example:"2023-05-15T00:00:00Z"`
	ReferalsAmount int    `json:"referals-amount" example:"10"`
	GamesAmount    int    `json:"games-amount" example:"5"`
}
