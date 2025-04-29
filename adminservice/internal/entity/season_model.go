package entity

import "time"

type Season struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	StartDate time.Time `json:"start-date"`
	EndDate   time.Time `json:"end-date"`
	Fund      uint      `json:"fund"`
}

type SeasonJson struct {
	StartDate string `json:"start-date" example:"2024-06-01" format:"date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end-date" example:"2024-08-31" format:"date" validate:"required,datetime=2006-01-02"`
}

type DetailSeasonJson struct {
	SeasonJson
	StartTime string `json:"start-time" example:"09:00:00" format:"time" validate:"required,datetime=15:04:05"`
	EndTime   string `json:"end-time" xample:"18:00:00" format:"time" validate:"required,datetime=15:04:05"`
	Fund      uint   `json:"fund" xample:"5000" minimum:"0" validate:"gte=0"`
}

type SeasonOut struct {
	ID        uint   `json:"id"`
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
	Fund      uint   `json:"fund"`
}
