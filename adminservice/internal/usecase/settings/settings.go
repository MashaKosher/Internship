package settings

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/di"
	"adminservice/internal/entity"
)

type UseCase struct {
	repo   repo.SettingsRepo
	logger di.LoggerType
	cfg    di.ConfigType
	bus    di.Bus
}

func New(r repo.SettingsRepo, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) *UseCase {
	return &UseCase{
		repo:   r,
		logger: logger,
		cfg:    cfg,
		bus:    bus,
	}
}

func (u *UseCase) UpdateSettings(settings entity.SettingsJson) error {
	dbSettiings := settings.ToDB()

	if err := u.repo.UpdateSettings(dbSettiings); err != nil {
		return err
	}

	go u.bus.GameSettingsProducer.SendGameSettings(settings)

	return nil
}

func (u *UseCase) GameSettings() (entity.SettingsJson, error) {

	settings, err := u.repo.GameSettings()
	if err != nil {
		return entity.SettingsJson{}, err
	}

	return settings.ToJSON(), nil
}
