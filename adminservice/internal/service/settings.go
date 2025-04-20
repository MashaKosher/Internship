package service

import (
	"adminservice/internal/entity"
	repo "adminservice/internal/repository"
	"adminservice/pkg"
	"encoding/json"
	"net/http"
	"time"
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

	answer, _ := r.Context().Value("val").(entity.AuthAnswer)
	if err := pkg.ValidateAuthResponse(answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var settings entity.SettingsJson
	var dbSettiings entity.GameSettings

	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbSettiings.LoseAmount = settings.LoseAmount
	dbSettiings.WinAmount = settings.WinAmount
	dbSettiings.WaitingTime = time.Duration(settings.WaitingTime)

	if err := repo.UpdateSettings(dbSettiings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// resStr := "The number of records is: " + string(amount)
	json.NewEncoder(w).Encode(settings)
}
