package settings

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"net/http"
)

type UseCase struct {
	repo repo.SettingsRepo
}

func New(r repo.SettingsRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (u *UseCase) UpdateSettings(w http.ResponseWriter, r *http.Request) (entity.SettingsJson, error) {
	if err := pkg.CheckToken(r); err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return entity.SettingsJson{}, err
	}

	var settings entity.SettingsJson
	var dbSettiings entity.GameSettings

	// Parsing Game Settings body
	if err := pkg.ParseGameSettingsBody(r.Body, &settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return entity.SettingsJson{}, err
	}

	// Filing Game Settings DB struct
	pkg.FillGameSettingsDBEntity(&settings, &dbSettiings)

	// Updating dbSettings in DB
	if err := u.repo.UpdateSettings(dbSettiings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return entity.SettingsJson{}, err
	}

	return settings, nil
}
