package service

import (
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"

	// "adminservice/internal/repository"
	repo "adminservice/internal/adapter/db/sql/plan"
	"adminservice/pkg"
	"encoding/json"
	"log"
	"net/http"
)

// @Summary Детальное палнирование сезона
// @Description Обрабатывает запрос на планирование сезона и проверяет права пользователя
// @Tags Season Planing
// @Accept json
// @Produce json
// @Param season body entity.DetailSeasonJson true "Информация о сезоне"
// @Success 200 {string} string "User is admin"
// @Router /deatil-plan [post]
func PlanSeason(w http.ResponseWriter, r *http.Request) {

	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var season entity.DetailSeasonJson
	var dbSeason entity.Season

	// parsing season body
	if err := pkg.ParseSeasonBody(r.Body, &season); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Serialize Season JSON entity in DB entity
	if err := pkg.StoreDeatailSeasonInDBEntity(&season, &dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Finding if seasons are crossing
	if err := repo.FindSeasonCross(&dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If season is not crossing with others we add it to DB
	if err := repo.AddNewSeason(&dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Season IN DB: ", dbSeason)

	// Produsing new season to Core service
	go producers.SendSeasonInfo(pkg.ParseSeasonToKafkaJSON(dbSeason))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dbSeason)
}
