package entity

type Season struct {
	ID        uint   `json:"id"`
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
	Fund      uint   `json:"fund"`
	Satatus   string `json:"status"`
}

type SeasonListElement struct {
	ID      uint   `json:"season-id"`
	Satatus string `json:"status"`
}

type Leaderboard struct {
	UserID int32 `json:"user-id"`
	Win    int32 `json:"win"`
}
