package service

import (
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"
	repo "adminservice/internal/repository"
	"adminservice/pkg"
	"encoding/json"
	"fmt"
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
	answer, _ := r.Context().Value("val").(entity.AuthAnswer)

	if err := pkg.ValidateAuthResponse(answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var DBDailyTasks entity.DBDailyTasks

	if err := json.NewDecoder(r.Body).Decode(&DBDailyTasks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DBDailyTasks.TaskDate = time.Now()

	fmt.Println("Текущая дата (без времени):", time.Now().Truncate(24*time.Hour))

	var DailyTasks entity.DailyTasks

	DailyTasks.GamesAmount = DBDailyTasks.GamesAmount
	DailyTasks.ReferalsAmount = DBDailyTasks.ReferalsAmount
	DailyTasks.TaskDate = DBDailyTasks.TaskDate.Format("2006-01-02")

	if err := repo.AddDailyTask(DBDailyTasks); err != nil {
		http.Error(w, "Daily Task for today already exists", http.StatusBadRequest)
		return
	}

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
	answer, _ := r.Context().Value("val").(entity.AuthAnswer)

	if err := pkg.ValidateAuthResponse(answer); err != nil {
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
