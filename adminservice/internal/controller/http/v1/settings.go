package v1

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"encoding/json"
	"net/http"
)

type settingsRoutes struct {
	u di.SettingsService
	v di.ValidatorType
	l di.LoggerType
}

func initSettingsRoutes(deps di.Container) *settingsRoutes {
	return &settingsRoutes{u: deps.Services.Settings, v: deps.Validator, l: deps.Logger}
}

// @Summary Update game settings
// @Description Update game configuration settings for authenticated user
// @Tags Settings
// @Accept  json
// @Produce  json
// @Param settings body entity.SettingsJson true "Game settings object"
// @Success 200 {object} entity.SettingsJson
// @Router			/settings [put]
func (gr *settingsRoutes) updatGameSettings(w http.ResponseWriter, r *http.Request) {
	var settings entity.SettingsJson

	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := gr.v.Struct(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := gr.u.UpdateSettings(settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(settings)
}

// @Summary Get game settings
// @Description Retrieve the game configuration settings for the current user
// @Tags Settings
// @Accept json
// @Produce json
// @Success 200 {object} entity.SettingsJson "Successfully retrieved game settings"
// @Failure 400 {object} entity.ErrorResponse "Game settings not found"
// @Failure 500 {object} entity.ErrorResponse "Internal server error"
// @Router /settings [get]
func (gr *settingsRoutes) gameSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := gr.u.GameSettings()
	if err != nil {
		if err == entity.ErrRecordNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(settings)
}
