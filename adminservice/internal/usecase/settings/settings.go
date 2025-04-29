package settings

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"adminservice/pkg"
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
	var dbSettiings entity.GameSettings

	// Filing Game Settings DB struct
	pkg.FillGameSettingsDBEntity(&settings, &dbSettiings)

	// Updating dbSettings in DB
	if err := u.repo.UpdateSettings(dbSettiings); err != nil {
		return err
	}

	go producers.SendGameSettings(settings, u.cfg, u.bus)

	return nil
}
