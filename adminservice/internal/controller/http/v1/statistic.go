package v1

import (
	"adminservice/internal/di"
	"net/http"
)

type statisticRoutes struct {
	u di.StatisticService
	v di.ValidatorType
	l di.LoggerType
}

func initStatisticRoutes(deps di.Container) *statisticRoutes {
	return &statisticRoutes{u: deps.Services.Statistic, v: deps.Validator, l: deps.Logger}
}

// @Summary		Get season statistic
// @Description	Get season statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/players [get]
func (sr *statisticRoutes) players(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	w.Write([]byte("There must be Players statistic"))
}

// @Summary		Get seasons statistic
// @Description	Get seasons statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/seasons [get]
func (sr *statisticRoutes) seasons(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	w.Write([]byte("There must be Season statistic"))
}
