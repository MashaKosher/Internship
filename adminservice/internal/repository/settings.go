package repository

import (
	"adminservice/internal/entity"
	db "adminservice/pkg/client/sql"
	"log"
)

func UpdateSettings(newSettings entity.GameSettings) error {

	var counter int64
	err := db.DB.Model(&entity.GameSettings{}).Count(&counter).Error
	if err != nil {
		return err
	}

	log.Println("Records amount: ", counter)

	if counter == 0 {
		if err := db.DB.Save(&newSettings).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Model(&entity.GameSettings{}).Where("id = 1").Updates(&newSettings).Error; err != nil {
			return err
		}
	}
	return err
}
