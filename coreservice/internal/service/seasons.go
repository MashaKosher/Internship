package service

import (
	"coreservice/internal/adapter/elastic"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"coreservice/internal/repository/sqlc"
	repo "coreservice/internal/repository/sqlc"
	"coreservice/pkg"
	"fmt"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
func SeasonInfo(c *gin.Context) {

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

	season, err := repo.GetSeasonById(int64(seasonID))
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "No such season"})
		return
	}

	logger.Logger.Info("ID:: " + paramID)
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
func Seasons(c *gin.Context) {

	seasons, err := repo.GetAllSeasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pkg.ConvertSeasonDBListToJson(seasons))

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
func SeasonLeaderBoard(c *gin.Context) {
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

	leaderBoard, err := repo.GetSeasonLeaderBoard(int64(seasonID))
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pkg.ConvertLeaderBoardDBListToJson(leaderBoard))
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
func CurrentSeason(c *gin.Context) {
	ids, err := elastic.SearchSeasonsByStatus(elastic.CurrentSeason)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	fmt.Println(ids)

	seasons, err := sqlc.GetSeasonsByIds(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

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
func PlannedSeason(c *gin.Context) {
	ids, err := elastic.SearchSeasonsByStatus(elastic.PlannedSeason)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	fmt.Println(ids)

	seasons, err := sqlc.GetSeasonsByIds(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, seasons)
}
