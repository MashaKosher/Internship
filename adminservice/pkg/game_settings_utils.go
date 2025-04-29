package pkg

import (
	"adminservice/internal/entity"
	"time"
)

func FillGameSettingsDBEntity(settings *entity.SettingsJson, dbSettiings *entity.GameSettings) {
	dbSettiings.LoseAmount = settings.LoseAmount
	dbSettiings.WinAmount = settings.WinAmount
	dbSettiings.WaitingTime = time.Duration(settings.WaitingTime)
}
