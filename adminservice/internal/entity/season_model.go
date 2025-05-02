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
	StartDate string `json:"start-date" example:"01-06-2024" format:"date" validate:"required,datetime=02-01-2006"`
	EndDate   string `json:"end-date" example:"31-08-3034" format:"date" validate:"required,datetime=02-01-2006"`
}

type DetailSeasonJson struct {
	SeasonJson
	StartTime string `json:"start-time" validate:"required"`
	EndTime   string `json:"end-time" validate:"required"`
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
