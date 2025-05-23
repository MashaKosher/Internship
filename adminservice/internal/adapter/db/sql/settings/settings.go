package settings

import (
	"adminservice/internal/entity"
	"errors"

	"gorm.io/gorm"
)

type SettingsRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *SettingsRepo {
	return &SettingsRepo{db}
}

func (r *SettingsRepo) UpdateSettings(newSettings *entity.GameSettings) error {
	if newSettings == nil {
		return entity.ErrGameSettingsIsNil
	}

	var counter int64
	err := r.DB.Model(&entity.GameSettings{}).Count(&counter).Error
	if err != nil {
		return err
	}

	if counter == 0 {
		if err := r.DB.Save(newSettings).Error; err != nil {
			return err
		}
	} else {
		if err := r.DB.Model(&entity.GameSettings{}).Where("id = 1").Updates(newSettings).Error; err != nil {
			return err
		}
	}
	return err
}

func (r *SettingsRepo) GameSettings() (entity.GameSettings, error) {
	var settings entity.GameSettings
	err := r.DB.First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.GameSettings{}, entity.ErrRecordNotFound
		}
		return entity.GameSettings{}, err
	}
	return settings, nil
}
