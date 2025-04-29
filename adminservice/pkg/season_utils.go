package pkg

import (
	"adminservice/internal/entity"
	"errors"
	"fmt"
	"log"
	"time"
)

func StoreDeatailSeasonInDBEntity(seasonJSON *entity.DetailSeasonJson, seasonDB *entity.Season) error {

	startTime, err := getTimeFromString(seasonJSON.StartDate, seasonJSON.StartTime)
	if err != nil {
		return err
	}

	if startTime.Before(time.Now()) {
		return errors.New("season cannot starting before Now")
	}

	endTime, err := getTimeFromString(seasonJSON.EndDate, seasonJSON.EndTime)
	if err != nil {
		return err
	}

	if startTime.After(endTime) {
		return errors.New("end date can't be earlier then start date")
	}

	seasonDB.StartDate = startTime
	seasonDB.EndDate = endTime
	seasonDB.Fund = seasonJSON.Fund

	return nil
}

func ParseSeasonToKafkaJSON(dbSeason entity.Season) entity.SeasonOut {
	var seasonOut entity.SeasonOut
	seasonOut.ID = dbSeason.ID
	seasonOut.StartDate = fmt.Sprint(dbSeason.StartDate)
	seasonOut.EndDate = fmt.Sprint(dbSeason.EndDate)
	seasonOut.Fund = dbSeason.Fund
	log.Println(seasonOut)
	return seasonOut
}
