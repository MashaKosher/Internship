package entity

type DailyTask struct {
	TaskDate       string `json:"task-date"`
	ReferalsAmount int    `json:"referals-amount"`
	GamesAmount    int    `json:"games-amount"`
}
