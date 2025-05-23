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

// @Summary Get today's daily task
// @Description Retrieve the daily task for the current date
// @Tags DailyTasks
// @Accept json
// @Produce json
// @Success 200 {object} entity.DailyTasks "Successfully retrieved today's task"
// @Failure 400 {object} entity.Response "No token"
// @Failure 401 {object} entity.Response "Invalid or expired token"
// @Failure 403 {object} entity.Response "User is not admin"
// @Failure 404 {object} entity.Response "Record not found in DB"
// @Failure 500 {object} entity.Response "Internal server error"
// @Router /daily-tasks [get]
func (dr *dailyTasksRoutes) dailyTask(w http.ResponseWriter, r *http.Request) {
	task, err := dr.u.GetDailyTask()
	if err != nil {
		if err == entity.ErrRecordNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// @Summary Create daily tasks
// @Description Create a new set of daily tasks for the current date
// @Tags DailyTasks
// @Accept json
// @Produce json
// @Param tasks body entity.DBDailyTasks true "Daily tasks data"
// @Example {json} Request-Example:
//
//		{
//		    "referals-amount": 10,
//	     "referals-task-reward": 15,
//		    "games-amount": 5,
//	     "game-task-reward": 10,
//		}
//
// @Success 201 {object} entity.DailyTasks "Successfully created today's task"
// @Failure 400 {object} entity.Response "No token or Invalid data"
// @Failure 401 {object} entity.Response "Invalid or expired token"
// @Failure 403 {object} entity.Response "User is not admin"
// @Failure 409 {object} entity.Response "Conflict: DailyTask already exists"
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
		if err == entity.ErrDailytaskIsNil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err == entity.ErrDailyTaskAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
// @Failure 400 {object} entity.Response "No token or Invalid data"
// @Failure 401 {object} entity.Response "Invalid or expired token"
// @Failure 403 {object} entity.Response "User is not admin"
// @Failure 404 {object} entity.Response "Record not found in DB"
// @Failure 500 {object} entity.Response "Internal server error"
// @Router /daily-tasks [delete]
func (dr *dailyTasksRoutes) deleteDailyTask(w http.ResponseWriter, r *http.Request) {

	err := dr.u.DeleteDailyTask()
	if err != nil {
		if err == entity.ErrRecordNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.Response{Message: "Task deleted successfully"})
}
