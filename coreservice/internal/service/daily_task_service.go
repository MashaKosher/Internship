package service

import (
	"net/http"

	repo "coreservice/internal/repository/sqlc"

	"github.com/gin-gonic/gin"
)

// @Summary Получить ежедневную задачу
// @Description Возвращает ежедневную задачу для текущей даты
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Router /daily-task [get]
func DailyTask(c *gin.Context) {

	task, err := repo.GetDailyTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}
