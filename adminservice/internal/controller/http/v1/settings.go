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
// @Router			/settings [post]
func (gr *settingsRoutes) gameSettings(w http.ResponseWriter, r *http.Request) {
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
