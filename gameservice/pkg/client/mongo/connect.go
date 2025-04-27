package mongo

import (
	"context"
	"fmt"
	"gameservice/internal/config"
	"gameservice/pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoConnect() (context.Context, *mongo.Client, *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%s", config.AppConfig.MongoDB.User, config.AppConfig.MongoDB.Password, config.AppConfig.MongoDB.Host, config.AppConfig.MongoDB.Port),
	))
	if err != nil {
		logger.L.Fatal(err.Error())
	}

	// Проверка подключения
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.L.Fatal(err.Error())
	}
	logger.L.Info("Successfully connected to MongoDB!")

	database := client.Database(config.AppConfig.MongoDB.Name)

	return ctx, client, database

}

// // 3. Создание коллекции (с опциями)
// collectionName := "users"
// opts := options.CreateCollection().SetValidator(bson.M{
// 	"$jsonSchema": bson.M{
// 		"bsonType": "object",
// 		"required": []string{"name", "email"},
// 		"properties": bson.M{
// 			"name": bson.M{
// 				"bsonType":    "string",
// 				"description": "must be a string and is required",
// 			},
// 			"email": bson.M{
// 				"bsonType":    "string",
// 				"pattern":     `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
// 				"description": "must be a valid email and is required",
// 			},
// 			"age": bson.M{
// 				"bsonType":    "int",
// 				"minimum":     18,
// 				"description": "must be an integer and >= 18",
// 			},
// 		},
// 	},
// })

// err = database.CreateCollection(ctx, collectionName, opts)
// if err != nil {
// 	// Если коллекция уже существует, получим ошибку, можно игнорировать
// 	// или обработать специальным образом
// 	fmt.Printf("Collection creation warning: %v\n", err)
// } else {
// 	fmt.Println("Collection created successfully!")
// }

// // 4. Получение ссылки на коллекцию для работы
// collection := database.Collection(collectionName)

// // 5. Пример вставки документа
// user := bson.M{
// 	"name":  "John Doe",
// 	"email": "john@example.com",
// 	"age":   30,
// }

// insertResult, err := collection.InsertOne(ctx, user)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedID)
