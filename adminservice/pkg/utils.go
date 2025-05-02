package pkg

import (
	"time"
)

func GetTimeFromString(stringDate, stringTime string) (time.Time, error) {
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

func parseDate(dateStr string) (time.Time, error) {
	layout := "02-01-2006"
	return time.Parse(layout, dateStr)
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
