package v1

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
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
	if r == nil || r.l == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	r.l.Info("We are in controller")

	_, exists := c.Get("message")
	if !exists {
		r.l.Error("Message not found in context")
		c.JSON(http.StatusBadRequest, entity.Error{Error: "message required"})
		return
	}
	r.l.Info("Message received")

	data, exists := c.Get("data")
	if !exists {
		r.l.Error("Data not found in context")
		c.JSON(http.StatusBadRequest, entity.Error{Error: "token data required"})
		return
	}
	r.l.Info("Data received")

	user, _ := data.(db.User)

	c.JSON(http.StatusOK, user)
}
