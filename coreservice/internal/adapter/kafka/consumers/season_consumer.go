package consumers

import (
	"coreservice/internal/config"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"coreservice/pkg"
	"fmt"
	"log"
	"time"

	"coreservice/internal/adapter/asynq/producer"
	"coreservice/internal/adapter/elastic"
	repo "coreservice/internal/repository/sqlc"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func RecieveSeasonInfo() {
	var err error
	var season entity.Season

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
		"group.id":          "authService",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		logger.Logger.Error("Failed to create consumer: " + err.Error())
	}
	defer consumer.Close()

	logger.Logger.Info("RecieveSeasonInfo is working")
	err = consumer.Subscribe(config.AppConfig.Kafka.SeasonTopicRecieve, nil)
	if err != nil {
		logger.Logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)

		if err == nil {

			log.Println("Maessage readed")
			logger.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeseriSeasonAnswer(msg.Value, season)
			if err != nil {
				logger.Logger.Error("Error while consuming: " + err.Error())
			}

			logger.Logger.Info(fmt.Sprint(answer))

			// add to db
			if err = repo.AddSeason(answer); err != nil {
				panic(err)
			}

			logger.Logger.Info("Start date: " + answer.StartDate + " End date: " + answer.EndDate)

			// Converting start season time
			// layout := "2006-01-02 15:04:05:00 +0300 +03"
			// startTime, _ := time.Parse(layout, answer.StartDate)

			// startTime, err := fixTimeFormat(answer.StartDate)

			// ds := time.Date(startTime.D)

			// fmt.Println()

			// correct := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), 0, time.UTC)

			// endTime, _ := time.Parse(layout, answer.EndDate)
			// // endTime, err := fixTimeFormat(answer.EndDate)

			// fmt.Println("\n", startTime, endTime, "\n")

			// fmt.Println(time.Now())

			// produce async tasks

			startTime, err := parseTimeToLocal(answer.StartDate)
			if err != nil {
				panic(err)
			}

			endTime, err := parseTimeToLocal(answer.EndDate)
			if err != nil {
				panic(err)
			}
			producer.PlanSeasonTasks(int(answer.ID), startTime, endTime)

			elastic.AddSeasonToIndex(int(answer.ID))
		} else {
			logger.Logger.Error("Error while consuming: " + err.Error())
		}
	}
}

func parseTimeToLocal(input string) (time.Time, error) {
	// Парсим строку в объект времени
	layout := "2006-01-02 15:04:05 -0700 MST"
	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		fmt.Println("Ошибка парсинга:", err)
		return time.Time{}, err
	}

	// Получаем локальную временную зону
	localLocation, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Ошибка загрузки временной зоны:", err)
		return time.Time{}, err
	}

	// Создаем новое время с теми же значениями, но в локальной временной зоне
	localTime := time.Date(
		parsedTime.Year(),
		parsedTime.Month(),
		parsedTime.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		parsedTime.Second(),
		parsedTime.Nanosecond(),
		localLocation,
	)

	// Выводим результат
	fmt.Println("Исходное время:", parsedTime)
	fmt.Println("Локальное время:", localTime)

	return localTime, nil
}

// func fixTimeFormat(timeStr string) (time.Time, error) {
// 	// Удаляем второй часовой пояс
// 	fixedStr := strings.Split(timeStr, " +0000")[0] + " +0000"

// 	// Парсим время
// 	t, err := time.Parse("2006-01-02 15:04:05 -0700", fixedStr)
// 	if err != nil {
// 		return time.Time{}, err
// 	}

// 	// Конвертируем в локальный часовой пояс
// 	return t.Local(), nil
// }
