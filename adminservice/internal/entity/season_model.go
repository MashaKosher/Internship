package entity

import "time"

type Season struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	StartDate time.Time `json:"start-date"`
	EndDate   time.Time `json:"end-date"`
	Fund      uint      `json:"fund"`
}

type SeasonJson struct {
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
}

type DetailSeasonJson struct {
	SeasonJson
	StartTime string `json:"start-time"`
	EndTime   string `json:"end-time"`
	Fund      uint   `json:"fund"`
}

type SeasonOut struct {
	ID        uint   `json:"id"`
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
	Fund      uint   `json:"fund"`
}
