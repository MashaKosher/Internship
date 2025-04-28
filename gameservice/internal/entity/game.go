package entity

import "time"

type GameResult struct {
	GameTime     time.Time
	Winner       int
	Loser        int
	WinAmount    float64
	LoseAmount   float64
	WinnerResult int
	LoserResult  int
}
