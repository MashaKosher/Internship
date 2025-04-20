package entity

import "time"

type DBDailyTasks struct {
	TaskDate       time.Time `json:"-" gorm:"primaryKey;type:date"`
	ReferalsAmount int       `json:"referals-amount"`
	GamesAmount    int       `json:"games-amount"`
}

type DailyTasks struct {
	TaskDate       string `json:"task-date"`
	ReferalsAmount int    `json:"referals-amount"`
	GamesAmount    int    `json:"games-amount"`
}
