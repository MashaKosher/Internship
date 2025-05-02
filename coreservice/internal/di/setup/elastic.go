package setup

import (
	"context"
	season "coreservice/internal/adapter/elastic/seasons"
	user "coreservice/internal/adapter/elastic/user"
	"coreservice/internal/di"
	"coreservice/internal/entity"

	userRepo "coreservice/internal/adapter/db/postgres/user"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func mustElastic(logger di.LoggerType, db di.DBType) di.ElasticType {
	client, err := elasticsearch.NewDefaultClient()

	if err != nil {
		panic(err)
	}

	logger.Info("Elastic connected successfully")

	go esCreateIndexIfNotexist(client, logger)

	seasonStatusRepo := season.New(client, entity.SeasonSearchIndex, logger)
	userNameRepo := user.New(client, entity.UserSearchIndex, logger, userRepo.New(db))

	return di.ElasticType{
		SeasonStatus: seasonStatusRepo,
		UserName:     userNameRepo,
	}
}

func esCreateIndexIfNotexist(ESClient di.ESClient, logger di.LoggerType) {

	createIndex(ESClient, logger, entity.UserSearchIndex)
	createIndex(ESClient, logger, entity.SeasonSearchIndex)

}

func createIndex(ESClient di.ESClient, logger di.LoggerType, indexType entity.ElasticIndexType) {
	_, err := esapi.IndicesExistsRequest{
		Index: []string{string(indexType)},
	}.Do(context.Background(), ESClient)

	// If error Index does not exists
	if err != nil {
		ESClient.Indices.Create(string(indexType))
	}

	logger.Info("Index " + string(indexType) + " Index exits")
}
