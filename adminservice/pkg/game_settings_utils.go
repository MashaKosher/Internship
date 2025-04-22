package pkg

import (
	"adminservice/internal/entity"
	"encoding/json"
	"io"
	"time"
)

func ParseGameSettingsBody(body io.ReadCloser, gameSetting *entity.SettingsJson) error {
	if err := json.NewDecoder(body).Decode(gameSetting); err != nil {
		return err
	}
	return nil
}

func FillGameSettingsDBEntity(settings *entity.SettingsJson, dbSettiings *entity.GameSettings) {
	dbSettiings.LoseAmount = settings.LoseAmount
	dbSettiings.WinAmount = settings.WinAmount
	dbSettiings.WaitingTime = time.Duration(settings.WaitingTime)
}
