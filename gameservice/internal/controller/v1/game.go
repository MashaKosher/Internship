package v1

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"gameservice/internal/controller/middlewares"
	"gameservice/internal/entity"
	"gameservice/internal/usecase"
	"gameservice/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type gameRoutes struct {
	u usecase.Game
}

func InitGameRoutes(e *echo.Echo, gameUseCase usecase.Game) {
	r := &gameRoutes{u: gameUseCase}

	group := e.Group("/game")
	group.Use(middlewares.CheckTokenMiddleWare())
	group.GET("/settings", r.gameSettings)
	group.GET("/check-token", r.checkToken)
	group.GET("/play", r.playGame)
	group.GET("/statistic", r.playerStatistic)
}

// GetGameSettings godoc
// @Summary Recieve Game Settings
// @Description Returns current  Game Settings (Win amount, Lose amount and Waitng Time)
// @Tags Game
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.GameSettings
// @Failure 400 {object} entity.Error
// @Router /game/settings [get]
func (r *gameRoutes) gameSettings(c echo.Context) error {
	gameSettings, err := r.u.GetGameSettings()
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, gameSettings)
}

// CheckToken godoc
// @Summary Validate JWT token
// @Description Verifies JWT token validity and returns user information
// @Tags Authentication
// @Produce json
// @Success 200 {object} entity.AuthAnswer "Token validation response with user data"
// @Failure 400 {object} entity.Error "Invalid token format"
// @Failure 401 {object} entity.Error "Missing or invalid token"
// @Router /game/check-token [get]
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
// @Router /game/statistic [get]
func (r *gameRoutes) playerStatistic(c echo.Context) error {
	data, exists := c.Get("data").(entity.AuthAnswer)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	playerStats, err := r.u.GetPlayerStatistic(int(data.ID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, playerStats)
}

// HTTP запрос для начала игры
// @Summary Start a new game
// @Description Start a game and wait for a second player to join.
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /game/play [get]
func (r *gameRoutes) playGame(c echo.Context) error {
	logger.L.Info("Starting Game")
	data, exists := c.Get("data").(entity.AuthAnswer)
	if !exists {
		return c.JSON(http.StatusBadRequest, entity.Error{Error: "some error"})
	}

	playerID := data.ID

	// Формируем URL для подключения к WebSocket
	u := url.URL{Scheme: "ws", Host: "localhost:8005", Path: "/ws", RawQuery: "player_id=" + strconv.Itoa(int(playerID))}
	logger.L.Info("Подключение к " + u.String())

	// Устанавливаем соединение
	wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.L.Error("Ошибка при подключении к WebSocket:" + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось подключиться к WebSocket"})
	}
	defer wsConn.Close()

	// Ожидаем сообщение от сервера
	_, message, err := wsConn.ReadMessage()
	if err != nil {
		logger.L.Error("Ошибка при чтении сообщения: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ошибка при чтении сообщения"})
	}

	// Декодируем JSON-ответ
	var result map[string]interface{}
	if err := json.Unmarshal(message, &result); err != nil {
		logger.L.Error("Ошибка при декодировании JSON:" + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "не удалось декодировать ответ"})
	}

	// Отправляем результат клиенту
	return c.JSON(http.StatusOK, result)
}
