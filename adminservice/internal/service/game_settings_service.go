package service

import (
	"adminservice/internal/entity"
	repo "adminservice/internal/repository"
	"adminservice/pkg"
	"encoding/json"
	"net/http"
)

// @Summary Update game settings
// @Description Update game configuration settings for authenticated user
// @Tags Settings
// @Accept  json
// @Produce  json
// @Param settings body entity.SettingsJson true "Game settings object"
// @Success 200 {object} entity.SettingsJson
// @Router			/settings [post]
func GameSettings(w http.ResponseWriter, r *http.Request) {

	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var settings entity.SettingsJson
	var dbSettiings entity.GameSettings

	// Parsing Game Settings body
	if err := pkg.ParseGameSettingsBody(r.Body, &settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Filing Game Settings DB struct
	pkg.FillGameSettingsDBEntity(&settings, &dbSettiings)

	// Updating dbSettings in DB
	if err := repo.UpdateSettings(dbSettiings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(settings)
}
