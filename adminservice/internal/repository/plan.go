package repository

import (
	"adminservice/internal/db"
	"adminservice/internal/entity"
	"errors"
	"log"
)

func AddNewSeason(season *entity.Season) error {
	if err := db.DB.Create(season).Error; err != nil {
		return err
	}
	return nil
}

func FindSeasonCross(season *entity.Season) error {
	var counter int64

	err := db.DB.Model(&entity.Season{}).Where("start_date <= ? AND end_date >= ?", season.StartDate, season.StartDate).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter > 0 {
		return errors.New("sesons are crossing")
	}

	err = db.DB.Model(&entity.Season{}).Where("start_date <= ? AND end_date >= ?", season.StartDate, season.EndDate).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter > 0 {
		return errors.New("sesons are crossing")
	}

	// Также проверяем, если новый сезон начинается после существующего
	err = db.DB.Model(&entity.Season{}).Where("start_date >= ? AND end_date <= ?", season.StartDate, season.EndDate).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter > 0 {
		return errors.New("sesons are crossing")
	}

	log.Println("Counter: ", counter)

	return nil
}
