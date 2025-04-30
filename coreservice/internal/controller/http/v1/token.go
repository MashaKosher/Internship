package v1

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokenRoutes struct {
	u di.TokenService
	l di.LoggerType
}

func TokenRoutes(app *gin.Engine, deps di.Container) {
	r := &tokenRoutes{l: deps.Logger}
	token := app.Group("/check-token")
	token.GET("/", r.checkToken)
}

// CheckToken godoc
// @Summary Validate JWT token
// @Description Verifies JWT token validity and returns user information
// @Tags Authentication
// @Produce json
// @Success 200 {object} entity.TypeResponse "Token validation response with user data"
// @Failure 400 {object} entity.Error "Invalid token format"
// @Failure 401 {object} entity.Error "Missing or invalid token"
// @Router /check-token [get]
func (r *tokenRoutes) checkToken(c *gin.Context) {
	message, exists := c.Get("message")
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
		return
	}

	data, exists := c.Get("data")
	if !exists {
		c.JSON(http.StatusBadRequest, entity.Error{Error: "Token is invalid"})
		return
	}

	resp, err := r.u.VerifyToken(message, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}
