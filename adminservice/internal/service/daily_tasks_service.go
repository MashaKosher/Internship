package service

import (
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"

	// repo "adminservice/internal/adapter/db/"

	repo "adminservice/internal/adapter/db/sql/daily_task"
	"adminservice/pkg"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// @Summary Create daily tasks
// @Description Create a new set of daily tasks
// @Tags DailyTasks
// @Accept  json
// @Produce  json
// @Param tasks body entity.DBDailyTasks true "Daily tasks object"
// @Success 201 {object} entity.DailyTasks
// @Router /daily-tasks [post]
func CreateDailyTasks(w http.ResponseWriter, r *http.Request) {

	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Adding task to DB
	var DBDailyTasks entity.DBDailyTasks
	if err := pkg.ParseDailyTaskBody(r.Body, &DBDailyTasks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	DBDailyTasks.TaskDate = time.Now()

	if err := repo.AddDailyTask(DBDailyTasks); err != nil {
		http.Error(w, "Daily Task for today already exists", http.StatusBadRequest)
		return
	}

	log.Println("Таска добавлена в БД успешно:")

	// Sending task to Core Service
	DailyTasks := pkg.ParseDailyTaskToKafkaJSON(DBDailyTasks)

	go producers.SendDailyTask(DailyTasks)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(DailyTasks)
}

// @Summary Delete daily tasks
// @Description Deelete todays task
// @Tags DailyTasks
// @Accept  json
// @Produce  json
// @Router /daily-tasks [delete]
func DeleteTodaysTask(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repo.DeleteTodaysTask(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Task deleted successfully")

}
