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
// @Failure 400 {object} entity.ErrorResponse "Invalid request format or validation error"
// @Failure 401 {object} entity.ErrorResponse "Unauthorized (invalid or missing token)"
// @Failure 403 {object} entity.ErrorResponse "Forbidden (user is not admin)"
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

	if err := sr.u.PlanSeasons(season); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(season)

}
