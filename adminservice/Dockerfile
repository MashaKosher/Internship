FROM golang:alpine

RUN go install github.com/air-verse/air@latest


RUN go install github.com/swaggo/swag/cmd/swag@latest


RUN apk add --no-cache librdkafka-dev


# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

RUN go get github.com/confluentinc/confluent-kafka-go/kafka

# Копирование остального кода
COPY ../ .

# Установка конфигурации air
COPY .air.toml ./

# Запуск air
CMD ["air"]