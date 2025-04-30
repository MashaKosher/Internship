package entity

import (
	"adminservice/pkg"
	"errors"
	"fmt"
	"time"
)

type Season struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	StartDate time.Time `json:"start-date"`
	EndDate   time.Time `json:"end-date"`
	Fund      uint      `json:"fund"`
}

func (s *Season) ToDTO() SeasonOut {
	return SeasonOut{
		ID:        s.ID,
		StartDate: fmt.Sprint(s.StartDate),
		EndDate:   fmt.Sprint(s.EndDate),
		Fund:      s.Fund,
	}
}

type SeasonJson struct {
	StartDate string `json:"start-date" example:"2024-06-01" format:"date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end-date" example:"2024-08-31" format:"date" validate:"required,datetime=2006-01-02"`
}

type DetailSeasonJson struct {
	SeasonJson
	StartTime string `json:"start-time" example:"09:00:00" format:"time" validate:"required,datetime=15:04:05"`
	EndTime   string `json:"end-time" example:"18:00:00" format:"time" validate:"required,datetime=15:04:05"`
	Fund      uint   `json:"fund" example:"5000" minimum:"0" validate:"gte=0"`
}

func (s *DetailSeasonJson) ToDB() (Season, error) {
	startTime, err := pkg.GetTimeFromString(s.StartDate, s.StartTime)
	if err != nil {
		return Season{}, err
	}

	if startTime.Before(time.Now()) {
		return Season{}, errors.New("season cannot starting before Now")
	}

	endTime, err := pkg.GetTimeFromString(s.EndDate, s.EndTime)
	if err != nil {
		return Season{}, err
	}

	if startTime.After(endTime) {
		return Season{}, errors.New("end date can't be earlier then start date")
	}

	return Season{StartDate: startTime, EndDate: endTime, Fund: s.Fund}, nil
}

type SeasonOut struct {
	ID        uint   `json:"id"`
	StartDate string `json:"start-date"`
	EndDate   string `json:"end-date"`
	Fund      uint   `json:"fund"`
}
