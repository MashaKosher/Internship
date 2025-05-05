package plan

import (
	"adminservice/internal/entity"
	"errors"
	"log"

	"gorm.io/gorm"
)

type PlanRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *PlanRepo {
	return &PlanRepo{db}
}

func (r *PlanRepo) AddNewSeason(season *entity.Season) error {
	if err := r.DB.Create(season).Error; err != nil {
		return err
	}
	return nil
}

func (r *PlanRepo) FindSeasonCross(season *entity.Season) error {
	var counter int64

	err := r.DB.Model(&entity.Season{}).Where("start_date <= ? AND end_date >= ?", season.StartDate, season.StartDate).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter > 0 {
		return errors.New("sesons are crossing")
	}

	err = r.DB.Model(&entity.Season{}).Where("start_date <= ? AND end_date >= ?", season.StartDate, season.EndDate).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter > 0 {
		return errors.New("sesons are crossing")
	}

	// Также проверяем, если новый сезон начинается после существующего
	err = r.DB.Model(&entity.Season{}).Where("start_date >= ? AND end_date <= ?", season.StartDate, season.EndDate).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter > 0 {
		return errors.New("sesons are crossing")
	}

	log.Println("Counter: ", counter)

	return nil
}

func (r *PlanRepo) Seasons() ([]entity.Season, error) {
	var seasons []entity.Season
	if err := r.DB.Find(&seasons).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Season{}, entity.ErrRecordNotFound
		}
		return []entity.Season{}, err
	}

	return seasons, nil
}
