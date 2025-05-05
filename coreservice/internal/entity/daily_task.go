package entity

type DailyTask struct {
	TaskDate           string `json:"task-date" example:"2023-05-15T00:00:00Z"`
	ReferalsAmount     int    `json:"referals-amount" example:"10"`
	ReferalsTaskReward int    `json:"referals-task-reward" example:"10"`
	GamesAmount        int    `json:"games-amount" example:"5"`
	GameTaskReward     int    `json:"game-task-reward" example:"10"`
}
