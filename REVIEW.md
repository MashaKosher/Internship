# Код-ревью: Сервисы Auth и Admin

## Ревью Auth сервиса

### Архитектура и дизайн

1. **Нарушения Clean Architecture**

   - В `internal/app/app.go` нарушен принцип инверсии зависимостей:

   ```go
   // Текущая реализация
   authUseCase := auth.New(
       authRepo.New(db),
   )
   ```

   Все зависимости создаются в одном месте, что делает код жестко связанным.

   - HTTP-контроллеры напрямую используют use cases без абстракций:

   ```go
   v1.NewRouter(app, authUseCase)
   ```

   Это нарушает принцип зависимости от абстракций.

2. **Проблемы с графом зависимостей**

   - Отсутствует контейнер зависимостей
   - Жесткая привязка к конкретным реализациям (Fiber, конкретные репозитории)
   - Сложная инициализация в `app.go`

3. **Структура проекта**

   Текущая структура:

   ```
   authservice/
   ├── cmd/
   │   └── auth/
   │       └── main.go
   ├── internal/
   │   ├── adapter/
   │   ├── controller/
   │   └── usecase/
   └── pkg/
   ```

   Рекомендуемая структура:

   ```
   authservice/
   ├── cmd/
   │   └── auth/
   │       └── main.go
   ├── internal/
   │   ├── domain/
   │   │   ├── entity/
   │   │   ├── repository/
   │   │   └── usecase/
   │   ├── infrastructure/
   │   │   ├── http/
   │   │   ├── kafka/
   │   │   └── db/
   │   └── di/
   │       └── container.go
   └── pkg/
   ```

### Проблемы и улучшения

1. **Обработка ошибок**

   - Непоследовательная обработка ошибок между слоями
   - Отсутствие пользовательских типов ошибок
   - Неправильное обертывание ошибок

   Рекомендация:

   ```go
   // Вместо
   if err != nil {
       return err
   }

   // Должно быть
   if err != nil {
       return fmt.Errorf("failed to process auth: %w", err)
   }
   ```

2. **Логирование**

   - Избыточное использование `log.Println`
   - Непоследовательные уровни логирования
   - Жестко закодированные сообщения логов

   Рекомендация:

   ```go
   // Вместо
   logger.Logger.Error("Server error:" + err.Error())

   // Должно быть
   logger.WithFields(logger.Fields{
       "error": err,
       "component": "server",
   }).Error("failed to start server")
   ```

3. **Безопасность**

   - Реализация хеширования паролей не видна в коде
   - Валидация токенов может быть более надежной
   - Отсутствие rate limiting

   Рекомендация:

   ```go
   // Добавить middleware для rate limiting
   app.Use(ratelimit.New(ratelimit.Config{
       Max:        100,
       Expiration: 1 * time.Minute,
   }))
   ```

4. **Качество кода**

   - Слишком длинные функции в `app.go`
   - Непоследовательное использование идиом Go
   - Использование "магических" строк и чисел

   Рекомендация:

   ```go
   // Вместо
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

   // Должно быть
   const shutdownTimeout = 10 * time.Second
   ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
   ```

## Ревью Admin сервиса

### Архитектура и дизайн

1. **Нарушения Clean Architecture**

   - Аналогичные проблемы с графом зависимостей
   - Смешивание бизнес-логики с HTTP-логикой
   - Отсутствие четких границ между слоями

2. **Проблемы с Kafka интеграцией**

   - Жестко закодированная конфигурация
   - Непоследовательная обработка ошибок
   - Отсутствие механизмов повторных попыток

   Рекомендация:

   ```go
   // Вместо
   producer, err := kafka.NewProducer(&kafka.ConfigMap{
       "bootstrap.servers": "localhost:9092",
   })

   // Должно быть
   type KafkaConfig struct {
       BootstrapServers string
       RetryCount      int
       RetryInterval   time.Duration
   }

   func NewProducer(cfg KafkaConfig) (*Producer, error) {
       // Инициализация с конфигурацией
   }
   ```

3. **Операции с базой данных**

   - Неоптимальные SQL-запросы
   - Отсутствие правильной обработки транзакций
   - Нет пулинга соединений

   Рекомендация:

   ```go
   // Вместо
   db.Exec("INSERT INTO users...")

   // Должно быть
   tx, err := db.Begin()
   if err != nil {
       return fmt.Errorf("failed to begin transaction: %w", err)
   }
   defer tx.Rollback()

   if _, err := tx.Exec("INSERT INTO users..."); err != nil {
       return fmt.Errorf("failed to insert user: %w", err)
   }

   if err := tx.Commit(); err != nil {
       return fmt.Errorf("failed to commit transaction: %w", err)
   }
   ```

## Общие рекомендации

1. **Конфигурация**ё

   - Реализовать правильную валидацию конфигурации
   - Чувствительную инфу хранить в енве
   - Вынести конфиги в директорию configs, где будут лежать yaml и т.п.

2. **Безопасность**

   - Реализовать правильную валидацию входных данных
   - Добавить ограничение частоты запросов
   - Реализовать правильное управление сессиями
   - Добавить заголовки безопасности

3. **Производительность**

   - Реализовать кэширование где это уместно
   - Оптимизировать запросы к базе данных
   - Добавить правильный пулинг соединений

4. **Мониторинг**

   - Добавить правильный сбор метрик
   - Реализовать проверки работоспособности
   - Добавить правильное логирование для мониторинга

## Заключение

Основные проблемы в обоих сервисах:

1. Нарушение принципов чистой архитектуры
2. Неправильное управление зависимостями
3. Смешивание слоев приложения
4. Отсутствие абстракций для внешних сервисов
5. Ты выносишь в pkg, очень сомнительные вещи. (pkg - код, который можно переиспользовать в другом проекте, например ты реализовал какую-нибудь полезную структуру данных или алгоритм, то это сюда)

Рекомендуется:

1. Реализовать правильный граф зависимостей через DI контейнер
2. Разделить слои приложения согласно Clean Architecture
3. Добавить абстракции для всех внешних сервисов
4. Улучшить обработку ошибок и логирование
5. Оптимизировать работу с базой данных и Kafka
