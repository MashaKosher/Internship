package pkg

import (
	"coreservice/internal/entity"
	db "coreservice/internal/repository/sqlc/generated"
	sqlcutils "coreservice/pkg/sqlc_utils"
	"errors"
	"fmt"
	"time"
)

func GetUserInfo(user *db.User) entity.User {
	return entity.User{
		Balance: sqlcutils.NumericToFloat64(user.Balance),
		WinRate: sqlcutils.NumericToFloat64(user.WinRate),
		Login:   user.Login,
	}
}

func ConvertDBUserSliceToUser(users []db.User) []entity.User {
	result := make([]entity.User, 0, len(users))
	for _, user := range users {
		result = append(result, GetUserInfo(&user))
	}
	return result
}

func ConvertAnyToDBUser(entity any) (db.User, error) {
	user, ok := entity.(db.User)
	if !ok {
		return db.User{}, errors.New("value cannot be converted to db.User")
	}

	return user, nil
}

func ParseTimeToLocal(input string) (time.Time, error) {
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
