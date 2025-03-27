FROM golang:alpine

RUN go install github.com/air-verse/air@latest


RUN go install github.com/swaggo/swag/cmd/swag@latest

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

# Копирование остального кода
COPY ../ .

# Установка конфигурации air
COPY .air.toml ./

# Запуск air
CMD ["air"]