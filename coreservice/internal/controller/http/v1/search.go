package v1

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchRoutes struct {
	u di.SearchService
	l di.LoggerType
	v di.ValidatorType
}

func SearchRoutes(router *gin.Engine, deps di.Container) {

	r := &searchRoutes{u: deps.Services.Search, l: deps.Logger, v: deps.Validator}

	search := router.Group("/elastic")
	search.GET("/create-index", r.userBuildSearchIndex)
	search.POST("/strict", r.searchElasticByNameStrict)
	search.POST("/wildcard", r.searchElasticByNameWildcard)
	search.POST("/fuzzy", r.searchElasticByNameFuzzy)
}

// UserBuildSearchIndex godoc
// @Summary Rebuild user search index
// @Description Recreates Elasticsearch index and imports all users from database
// @Tags Elastic
// @Produce json
// @Success 200 {object} entity.Response "Returns count of indexed users"
// @Failure 500 {object} entity.Error "Elasticsearch operation failed"
// @Router /elastic/create-index [get]
func (r *searchRoutes) userBuildSearchIndex(c *gin.Context) {
	users, err := r.u.UserBuildSearchIndex()
	if err != nil {
		c.JSON(http.StatusOK, entity.Error{Error: err.Error()})
		return
	}
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
func (r *searchRoutes) searchElasticByNameStrict(c *gin.Context) {
	var searchParams entity.SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	if err := r.v.Struct(searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Invalid Search Param"})
		return
	}

	users, err := r.u.SearchElasticByNameStrict(searchParams.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
		return
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
func (r *searchRoutes) searchElasticByNameWildcard(c *gin.Context) {
	var searchParams entity.SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	if err := r.v.Struct(searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Invalid Search Param"})
		return
	}

	users, err := r.u.SearchElasticByNameWildcard(searchParams.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
		return
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
func (r *searchRoutes) searchElasticByNameFuzzy(c *gin.Context) {
	var searchParams entity.SearchParams
	if err := c.ShouldBindJSON(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	if err := r.v.Struct(searchParams); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Invalid Search Param"})
		return
	}

	users, err := r.u.SearchElasticByNameFuzzy(searchParams.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
