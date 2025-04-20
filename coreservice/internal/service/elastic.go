package service

import (
	"coreservice/internal/adapter/elastic"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserBuildSearchIndex godoc
// @Summary Rebuild user search index
// @Description Recreates Elasticsearch index and imports all users from database
// @Tags Elastic
// @Produce json
// @Success 200 {object} entity.Response "Returns count of indexed users"
// @Failure 500 {object} entity.Error "Elasticsearch operation failed"
// @Router /elastic/create-index [get]
func UserBuildSearchIndex(c *gin.Context) {
	users, err := elastic.AddingAllUsersToIndex()
	if err != nil {
		c.JSON(http.StatusOK, entity.Error{Error: err.Error()})
	}
	logger.Logger.Info("Users from DB: " + fmt.Sprint(users))
	c.JSON(http.StatusOK, users)

}

// SearchElasticByNameStrict godoc
// @Summary Search users by exact name
// @Description Performs a case-sensitive exact match search for users by name
// @Tags Elastic
// @Accept json
// @Produce json
// @Param request body entity.SearchParams true "Search parameters"
// @Success 200 {array} entity.User "List of found users"
// @Failure 400 {object} entity.Error  "Invalid request format"
// @Failure 500 {object} entity.Error  "Elasticsearch error"
// @Router /elastic/strict [post]
func SearchElasticByNameStrict(c *gin.Context) {
	var searchParams entity.SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	users, err := elastic.GetUserByName(searchParams.Name, elastic.Strict)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
	}
	c.JSON(http.StatusOK, users)
}

// SearchElasticByNameWildcard godoc
// @Summary Search users by name using wildcard pattern
// @Description Performs a wildcard search for users by name (supports * and ? patterns)
// @Tags Elastic
// @Accept json
// @Produce json
// @Param request body entity.SearchParams true "Search parameters"
// @Success 200 {array} entity.User "List of found users"
// @Failure 400 {object} entity.Error  "Invalid request format"
// @Failure 500 {object} entity.Error  "Elasticsearch error"
// @Router /elastic/wildcard [post]
func SearchElasticByNameWildcard(c *gin.Context) {
	var searchParams entity.SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	users, err := elastic.GetUserByName(searchParams.Name, elastic.Wildcard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
	}
	c.JSON(http.StatusOK, users)
}

// SearchElasticByNameFuzzy godoc
// @Summary Fuzzy search users by name
// @Description Performs a fuzzy search that accounts for typos and similar spellings
// @Tags Elastic
// @Accept json
// @Produce json
// @Param request body entity.SearchParams true "Search parameters"
// @Success 200 {array} entity.User "List of found users"
// @Failure 400 {object} entity.Error  "Invalid request format"
// @Failure 500 {object} entity.Error  "Elasticsearch error"
// @Router /elastic/fuzzy [post]
func SearchElasticByNameFuzzy(c *gin.Context) {
	var searchParams entity.SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	users, err := elastic.GetUserByName(searchParams.Name, elastic.Fuzzy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
	}
	c.JSON(http.StatusOK, users)
}
