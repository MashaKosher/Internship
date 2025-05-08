package v1

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"encoding/json"
	"errors"
	"net/http"
)

type seasonRoutes struct {
	u di.SesaonPlanService
	v di.ValidatorType
	l di.LoggerType
}

func initPlanRoutes(deps di.Container) *seasonRoutes {
	return &seasonRoutes{u: deps.Services.Plan, v: deps.Validator, l: deps.Logger}
}

// @Summary Season detailed planning
// @Description Handles season planning request and verifies user admin rights
// @Tags Season Planning
// @Accept json
// @Produce json
// @Param season body entity.DetailSeasonJson true "Season information"
// @Example {json} Request Example:
//
//	{
//	    "start-date": "2024-06-01",
//	    "end-date": "2024-08-31",
//	    "start-time": "09-00-00",
//	    "end-time": "18-00-00",
//	    "fund": 5000
//	}
//
// @Success 201 {object} entity.DetailSeasonJson "Successfully created season plan"
// @Failure 400 {object} entity.Response "No token or Invalid data"
// @Failure 401 {object} entity.Response "Invalid or expired token"
// @Failure 403 {object} entity.Response "User is not admin"
// @Failure 409 {object} entity.Response "Conflict: Seasons are crossing"
// @Failure 500 {object} entity.ErrorResponse "Internal server error"
// @Router /deatil-plan [post]
func (sr *seasonRoutes) planSeason(w http.ResponseWriter, r *http.Request) {
	var season entity.DetailSeasonJson

	if err := json.NewDecoder(r.Body).Decode(&season); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := sr.v.Struct(season)
	if err != nil {
		http.Error(w, errors.New("invalid json format").Error(), http.StatusBadRequest)
		return
	}

	if err := sr.u.PlanSeason(season); err != nil {
		if err == entity.ErrSeasonIsNil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err == entity.ErrSeasonsAreCrossing {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(season)
}

// @Summary Get all seasons
// @Description Retrieve the list of all available seasons
// @Tags Season Planning
// @Accept json
// @Produce json
// @Success 200 {array} entity.Season "Successfully retrieved all seasons"
// @Failure 400 {object} entity.Response "No token or Invalid data"
// @Failure 401 {object} entity.Response "Invalid or expired token"
// @Failure 403 {object} entity.Response "User is not admin"
// @Failure 500 {object} entity.ErrorResponse "Internal server error"
// @Router /seasons [get]
func (sr *seasonRoutes) seasons(w http.ResponseWriter, r *http.Request) {
	seasons, err := sr.u.Seasons()
	if err != nil {
		if err == entity.ErrRecordNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	prettyJSON, err := json.MarshalIndent(seasons, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)

}
