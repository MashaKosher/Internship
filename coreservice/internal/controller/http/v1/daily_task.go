package v1

import (
	"coreservice/internal/di"
	"net/http"

	"github.com/gin-gonic/gin"
)

type dailyTasksRoutes struct {
	u di.DailyTaskService
	l di.LoggerType
}

func DailyTasksRoutes(router *gin.Engine, deps di.Container) {
	r := &dailyTasksRoutes{u: deps.Services.DailyTask, l: deps.Logger}
	router.GET("/daily-task", r.dailyTask)

}

// @Summary Получить ежедневную задачу
// @Description Возвращает ежедневную задачу для текущей даты
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Router /daily-task [get]
func (r *dailyTasksRoutes) dailyTask(c *gin.Context) {
	task, err := r.u.DailyTask()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, task)
}
