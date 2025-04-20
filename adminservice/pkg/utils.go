package pkg

import (
	"adminservice/internal/entity"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

func ValidateAuthResponse(answer entity.AuthAnswer) error {
	if len(answer.Err) != 0 {
		return errors.New(answer.Err)
	}

	if answer.Role != "admin" {
		return errors.New("user is not admin")
	}

	return nil
}

func ParseResponse(body io.ReadCloser, season *entity.SeasonJson) error {
	if err := json.NewDecoder(body).Decode(season); err != nil {
		return err
	}
	return nil
}

func ParseDetailResponse(body io.ReadCloser, season *entity.DetailSeasonJson) error {
	if err := json.NewDecoder(body).Decode(season); err != nil {
		return err
	}
	return nil
}

func parseDate(dateStr string) (time.Time, error) {
	layout := "02-01-2006"
	return time.Parse(layout, dateStr)
}

func StoreSeasonInDBEntity(seasonJSON *entity.SeasonJson, seasonDB *entity.Season) error {

	startDate, err := parseDate(seasonJSON.StartDate)
	if err != nil {
		return err
	}

	log.Println("Start date before now " + fmt.Sprint(startDate.Before(time.Now())))
	log.Println("Start date: " + fmt.Sprint(startDate))
	log.Println("Now: " + fmt.Sprint(time.Now()))

	if startDate.Before(time.Now()) {
		return errors.New("sesons cannot starting before Now")
	}

	endDate, err := parseDate(seasonJSON.EndDate)
	if err != nil {
		return err
	}

	if startDate.After(endDate) {
		return errors.New("end date can't be earlier then start date")
	}

	log.Println("Dates: " + fmt.Sprintln(startDate.After(endDate)))

	seasonDB.StartDate = startDate
	seasonDB.EndDate = endDate

	return nil
}

func parseTime(timeStr string) (time.Time, error) {
	layout := "15-04-05"
	return time.Parse(layout, timeStr)
}

func StoreDeatailSeasonInDBEntity(seasonJSON *entity.DetailSeasonJson, seasonDB *entity.Season) error {

	startDate, err := parseDate(seasonJSON.StartDate)
	if err != nil {
		return err
	}

	startTime, err := parseTime(seasonJSON.StartTime)
	if err != nil {
		return err
	}

	start := time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(),
		0,
		time.UTC,
	)

	// log.Println("Start date before now " + fmt.Sprint(startDate.Before(time.Now())))
	// log.Println("Start date: " + fmt.Sprint(startDate))
	// log.Println("Now: " + fmt.Sprint(time.Now()))

	if start.Before(time.Now()) {
		return errors.New("sesons cannot starting before Now")
	}

	endDate, err := parseDate(seasonJSON.StartDate)
	if err != nil {
		return err
	}

	endTime, err := parseTime(seasonJSON.EndTime)
	if err != nil {
		return err
	}

	end := time.Date(
		endDate.Year(),
		endDate.Month(),
		endDate.Day(),
		endTime.Hour(),
		endTime.Minute(),
		endTime.Second(),
		0,
		time.UTC,
	)

	if start.After(end) {
		return errors.New("end date can't be earlier then start date")
	}

	// log.Println("Dates: " + fmt.Sprintln(startDate.After(endDate)))

	seasonDB.StartDate = start
	seasonDB.EndDate = end
	seasonDB.Fund = seasonJSON.Fund

	return nil
}
