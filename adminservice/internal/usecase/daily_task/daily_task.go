package dailytask

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"log"
	"net/http"
	"time"
)

type UseCase struct {
	// Пиши интерфейсы по месту использования, а не реализации.
	// Интерфейсы - контракт, заключаемый между вызывающим и вызываемым кодом.
	repo repo.DailyTaskRepo
}

func New(r repo.DailyTaskRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) CreateDailyTask(w http.ResponseWriter, r *http.Request) (entity.DailyTasks, error) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		return entity.DailyTasks{}, err
	}

	// Adding task to DB
	var DBDailyTasks entity.DBDailyTasks
	if err := pkg.ParseDailyTaskBody(r.Body, &DBDailyTasks); err != nil {
		return entity.DailyTasks{}, err
	}
	DBDailyTasks.TaskDate = time.Now()

	if err := uc.repo.AddDailyTask(DBDailyTasks); err != nil {
		return entity.DailyTasks{}, err
	}

	log.Println("Таска добавлена в БД успешно:")

	// Sending task to Core Service
	DailyTasks := pkg.ParseDailyTaskToKafkaJSON(DBDailyTasks)

	go producers.SendDailyTask(DailyTasks)

	return DailyTasks, nil
}

func (uc *UseCase) DeleteDailyTask(w http.ResponseWriter, r *http.Request) error {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := uc.repo.DeleteTodaysTask(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}
