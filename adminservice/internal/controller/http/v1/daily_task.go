package v1

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type dailyTasksRoutes struct {
	u di.DailyTaskService
	v di.ValidatorType
	l di.LoggerType
}

func initDailyTaskRoutes(deps di.Container) *dailyTasksRoutes {
	return &dailyTasksRoutes{u: deps.Services.DailyTask, v: deps.Validator, l: deps.Logger}
}

// @Summary Create daily tasks
// @Description Create a new set of daily tasks for the current date
// @Tags DailyTasks
// @Accept json
// @Produce json
// @Param tasks body entity.DBDailyTasks true "Daily tasks data"
// @Example {json} Request-Example:
//
//	{
//	    "referals-amount": 10,
//	    "games-amount": 5
//	}
//
// @Success 201 {object} entity.DailyTasks "Successfully created daily tasks"
// @Failure 400 {object} entity.Response "Invalid request format or validation error"
// @Failure 500 {object} entity.Response "Internal server error"
// @Router /daily-tasks [post]
func (dr *dailyTasksRoutes) createDailyTask(w http.ResponseWriter, r *http.Request) {
	var dailyTask entity.DBDailyTasks

	if err := json.NewDecoder(r.Body).Decode(&dailyTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dailyTask.TaskDate = time.Now()

	err := dr.v.Struct(dailyTask)
	if err != nil {
		http.Error(w, errors.New("invalid json format").Error(), http.StatusBadRequest)
		return
	}

	dailyTaskOut, err := dr.u.CreateDailyTask(dailyTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dailyTaskOut)
}

// @Summary Delete today's daily tasks
// @Description Delete daily tasks for the current date
// @Tags DailyTasks
// @Accept json
// @Produce json
// @Success 200 {object} entity.Response "Successfully deleted today's tasks"
// @Failure 400 {object} entity.Response "No tasks found for today or deletion error"
// @Failure 500 {object} entity.Response "Internal server error"
// @Router /daily-tasks [delete]
func (dr *dailyTasksRoutes) deleteTodaysTask(w http.ResponseWriter, r *http.Request) {

	err := dr.u.DeleteDailyTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.Response{Message: "Task deleted successfully"})
}
