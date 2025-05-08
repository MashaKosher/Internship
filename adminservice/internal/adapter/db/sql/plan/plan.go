package plan

import (
	"adminservice/internal/entity"
	"errors"

	"gorm.io/gorm"
)

type PlanRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *PlanRepo {
	return &PlanRepo{db}
}

func (r *PlanRepo) AddNewSeason(season *entity.Season) error {
	if season == nil {
		return entity.ErrSeasonIsNil
	}

	var count int64
	err := r.DB.Model(&entity.Season{}).
		Where("start_date <= ? AND end_date >= ?", season.EndDate, season.StartDate).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return entity.ErrSeasonsAreCrossing
	}

	return r.DB.Create(season).Error
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
