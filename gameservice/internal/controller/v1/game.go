package v1

import (
	"net/http"

	"gameservice/internal/controller/middlewares"
	"gameservice/internal/entity"
	"gameservice/internal/usecase"

	"github.com/labstack/echo/v4"
)

type gameRoutes struct {
	u usecase.Game
}

func InitGameRoutes(e *echo.Echo, gameUseCase usecase.Game) {
	r := &gameRoutes{u: gameUseCase}

	e.GET("/check-token", r.checkToken)

	group := e.Group("/game")
	group.Use(middlewares.CheckTokenMiddleWare())
	group.GET("/play", r.playGame)
	group.GET("/statistic", r.playerStatistic)
}

// CheckToken godoc
// @Summary Validate JWT token
// @Description Verifies JWT token validity and returns user information
// @Tags Authentication
// @Produce json
// @Success 200 {object} entity.AuthAnswer "Token validation response with user data"
// @Failure 400 {object} entity.Error "Invalid token format"
// @Failure 401 {object} entity.Error "Missing or invalid token"
// @Router /check-token [get]
func (r *gameRoutes) checkToken(c echo.Context) error {
	// Получаем данные из контекста
	_, exists := c.Get("message").(string)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	data, exists := c.Get("data").(entity.AuthAnswer)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	return c.JSON(http.StatusOK, data)
}

// @Tags Game
// @Produce json
// @Router /game/play [get]
func (r *gameRoutes) playGame(c echo.Context) error {
	// Получаем данные из контекста
	_, exists := c.Get("message").(string)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	data, exists := c.Get("data").(entity.AuthAnswer)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	return c.JSON(http.StatusOK, data)
}

// @Tags Game
// @Produce json
// @Router /game/statistic [get]
func (r *gameRoutes) playerStatistic(c echo.Context) error {
	// Получаем данные из контекста
	_, exists := c.Get("message").(string)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	data, exists := c.Get("data").(entity.AuthAnswer)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	return c.JSON(http.StatusOK, data)
}
