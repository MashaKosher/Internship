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

func ParseResponse(body io.ReadCloser, season *entity.SeasonJson) error {
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

func joinTime(joinDate time.Time, joinTime time.Time) time.Time {
	return time.Date(
		joinDate.Year(),
		joinDate.Month(),
		joinDate.Day(),
		joinTime.Hour(),
		joinTime.Minute(),
		joinTime.Second(),
		0,
		time.UTC,
	)
}

func getTimeFromString(stringDate, stringTime string) (time.Time, error) {
	timeFormatDate, err := parseDate(stringDate)
	if err != nil {
		return time.Time{}, err
	}

	timeFormatTime, err := parseTime(stringTime)
	if err != nil {
		return time.Time{}, err
	}

	mergedTime := joinTime(timeFormatDate, timeFormatTime)
	return mergedTime, nil
}
