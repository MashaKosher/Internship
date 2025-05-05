package entity

import "time"

type DBDailyTasks struct {
	TaskDate           time.Time `json:"-" gorm:"primaryKey;type:date" validate:"required"`
	ReferalsAmount     int       `json:"referals-amount" example:"10" minimum:"0" validate:"gte=0"`
	ReferalsTaskReward int       `json:"referals-task-reward" example:"10" minimum:"0" validate:"gte=0"`
	GamesAmount        int       `json:"games-amount" example:"5" minimum:"0" validate:"gte=0"`
	GameTaskReward     int       `json:"game-task-reward" example:"10" minimum:"0" validate:"gte=0"`
}

func (task *DBDailyTasks) ToDTO() DailyTasks {
	return DailyTasks{
		TaskDate:           task.TaskDate.Format("2006-01-02"),
		ReferalsAmount:     task.ReferalsAmount,
		ReferalsTaskReward: task.ReferalsTaskReward,
		GamesAmount:        task.GamesAmount,
		GameTaskReward:     task.GameTaskReward,
	}
}

type DailyTasks struct {
	TaskDate           string `json:"task-date" example:"2023-05-15T00:00:00Z"`
	ReferalsAmount     int    `json:"referals-amount" example:"10"`
	ReferalsTaskReward int    `json:"referals-task-reward" example:"10"`
	GamesAmount        int    `json:"games-amount" example:"5"`
	GameTaskReward     int    `json:"game-task-reward" example:"10"`
}
