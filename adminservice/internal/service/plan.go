package service

import (
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"
	"adminservice/internal/repository"
	"adminservice/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// @Summary Планирование сезона
// @Description Обрабатывает запрос на планирование сезона и проверяет права пользователя
// @Tags Season Planing
// @Accept json
// @Produce json
// @Param season body entity.SeasonJson true "Информация о сезоне"
// @Success 200 {string} string "User is admin"
// @Router /plan [post]
func PlanSeason(w http.ResponseWriter, r *http.Request) {
	answer, _ := r.Context().Value("val").(entity.AuthAnswer)

	if err := pkg.ValidateAuthResponse(answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var season entity.SeasonJson
	var dbSeason entity.Season

	if err := pkg.ParseResponse(r.Body, &season); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pkg.StoreSeasonInDBEntity(&season, &dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.FindSeasonCross(&dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.AddNewSeason(&dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(season)
	log.Println(dbSeason)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dbSeason)

}

// @Summary Детальное палнирование сезона
// @Description Обрабатывает запрос на планирование сезона и проверяет права пользователя
// @Tags Season Planing
// @Accept json
// @Produce json
// @Param season body entity.DetailSeasonJson true "Информация о сезоне"
// @Success 200 {string} string "User is admin"
// @Router /deatil-plan [post]
func DetilPlanSeason(w http.ResponseWriter, r *http.Request) {
	answer, _ := r.Context().Value("val").(entity.AuthAnswer)

	if err := pkg.ValidateAuthResponse(answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var season entity.DetailSeasonJson
	var dbSeason entity.Season

	if err := pkg.ParseDetailResponse(r.Body, &season); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := pkg.StoreDeatailSeasonInDBEntity(&season, &dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.FindSeasonCross(&dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.AddNewSeason(&dbSeason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(season)
	log.Println(dbSeason)

	w.WriteHeader(http.StatusCreated)

	log.Println("Season IN DB: ", dbSeason)

	var seasonOut entity.SeasonOut
	seasonOut.ID = dbSeason.ID
	seasonOut.StartDate = fmt.Sprint(dbSeason.StartDate)
	seasonOut.EndDate = fmt.Sprint(dbSeason.EndDate)
	seasonOut.Fund = dbSeason.Fund

	log.Println(seasonOut)

	go producers.SendSeasonInfo(seasonOut)
	json.NewEncoder(w).Encode(dbSeason)
}
