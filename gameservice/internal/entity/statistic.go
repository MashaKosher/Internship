package entity

type PlayerStats struct {
	PlayerID    int `json:"player-id"`
	TotalGames  int `json:"total-games"`
	TotalWins   int `json:"total-wins"`
	TotalLosses int `json:"total-losses"`
}
