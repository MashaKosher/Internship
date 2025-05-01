package v1

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type seasonRoutes struct {
	u di.SeasonService
	l di.LoggerType
	v di.ValidatorType
}

func SeasonRoutes(router *gin.Engine, deps di.Container) {

	r := &seasonRoutes{u: deps.Services.Season, l: deps.Logger, v: deps.Validator}
	season := router.Group("/seasons")
	season.GET("/", r.seasons)
	season.GET("/current", r.currentSeason)
	season.GET("/planned", r.plannedSeason)
	season.GET("/:id/leader-board", r.seasonLeaderBoard)
	season.GET("/:id", r.seasonInfo)
}

// @Summary Get season by ID
// @Description Get season information by season ID
// @Tags Seasons
// @Accept  json
// @Produce  json
// @Param id path int true "Season ID"
// @Success 200 {object} entity.Season
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Router /seasons/{id} [get]
func (r *seasonRoutes) seasonInfo(c *gin.Context) {
	paramID := c.Param("id")
	if len(paramID) == 0 {
		c.JSON(http.StatusBadRequest, "Invalid param")
		return
	}

	seasonID, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Invlaid param value"})
		return
	}

	season, err := r.u.SeasonInfo(seasonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "No such season"})
		return
	}
	r.l.Info("Sesaon Info for Season with ID: " + paramID + " successfully found")
	c.JSON(http.StatusOK, season)
}

// @Summary Get all seasons
// @Description Get a list of all seasons
// @Tags Seasons
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.SeasonListElement
// @Failure 500 {object} entity.Error
// @Router /seasons/ [get]
func (r *seasonRoutes) seasons(c *gin.Context) {
	seasons, err := r.u.Seasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, seasons)

}

// @Summary Get season leaderboard
// @Description Get leaderboard for specific season by ID
// @Tags Seasons
// @Accept  json
// @Produce  json
// @Param id path int true "Season ID"
// @Success 200 {array} entity.Leaderboard
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /seasons/{id}/leader-board [get]
func (r *seasonRoutes) seasonLeaderBoard(c *gin.Context) {
	paramID := c.Param("id")
	if len(paramID) == 0 {
		c.JSON(http.StatusBadRequest, "Invlaid param")
		return
	}
	seasonID, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Invlaid param value"})
		return
	}

	leaderBoard, err := r.u.SeasonLeaderBoard(seasonID)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}
	r.l.Info("Leader Board for Season with ID: " + paramID + " successfully found")
	c.JSON(http.StatusOK, leaderBoard)
}

// CurrentSeason godoc
// @Summary Get current seasons
// @Description Returns a list of all seasons with 'current' status
// @Tags Seasons
// @Accept  json
// @Produce json
// @Failure 400 {object} entity.Error "Bad request"
// @Failure 500 {object} entity.Error "Internal server error"
// @Router /seasons/current [get]
func (r *seasonRoutes) currentSeason(c *gin.Context) {

	seasons, err := r.u.CurrentSeason()
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	r.l.Info("Current season succesfully found: " + fmt.Sprint(seasons))
	c.JSON(http.StatusOK, seasons)
}

// PlannedSeason godoc
// @Summary Get planned seasons
// @Description Returns a list of all seasons with 'planned' status
// @Tags Seasons
// @Accept  json
// @Produce json
// @Failure 400 {object} entity.Error "Bad request"
// @Failure 500 {object} entity.Error "Internal server error"
// @Router /seasons/planned [get]
func (r *seasonRoutes) plannedSeason(c *gin.Context) {
	seasons, err := r.u.PlannedSeason()
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}
	r.l.Info("Planned seasons succesfully: " + fmt.Sprint(seasons))
	c.JSON(http.StatusOK, seasons)
}
