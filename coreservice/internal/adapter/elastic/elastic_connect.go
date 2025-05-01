package elastic

// import (
// 	"context"
// 	"coreservice/internal/logger"

// 	"github.com/elastic/go-elasticsearch/v8"
// 	"github.com/elastic/go-elasticsearch/v8/esapi"
// )

// var ESClient *elasticsearch.Client
// var UserSearchIndex string = "users"
// var SeasonSearchIndex string = "seasons"

// func ESClientConnection() {
// 	client, err := elasticsearch.NewDefaultClient()

// 	if err != nil {
// 		panic(err)
// 	}

// 	logger.Logger.Info("Elastic connected successfully")

// 	ESClient = client
// }

// func ESCreateIndexIfNotexist() {
// 	_, err := esapi.IndicesExistsRequest{
// 		Index: []string{UserSearchIndex},
// 	}.Do(context.Background(), ESClient)

// 	// If error Index does not exists
// 	if err != nil {
// 		ESClient.Indices.Create(UserSearchIndex)
// 	}

// 	logger.Logger.Info("Index User Index exitss")

// 	_, err = esapi.IndicesExistsRequest{
// 		Index: []string{SeasonSearchIndex},
// 	}.Do(context.Background(), ESClient)

// 	// If error Index does not exists
// 	if err != nil {
// 		ESClient.Indices.Create(SeasonSearchIndex)
// 	}

// 	logger.Logger.Info("Index Season Index exitss")

// }
