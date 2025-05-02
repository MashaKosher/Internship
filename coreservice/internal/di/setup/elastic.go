package setup

import (
	"context"
	"coreservice/internal/di"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const UserSearchIndex string = "users"
const SeasonSearchIndex string = "seasons"

func mustElastic(logger di.LoggerType) di.ElasticType {
	client, err := elasticsearch.NewDefaultClient()

	if err != nil {
		panic(err)
	}

	logger.Info("Elastic connected successfully")

	go esCreateIndexIfNotexist(client, logger)

	return di.ElasticType{
		ESClient:          client,
		UserSearchIndex:   UserSearchIndex,
		SeasonSearchIndex: SeasonSearchIndex,
	}
}

func esCreateIndexIfNotexist(ESClient di.ESClient, logger di.LoggerType) {

	createIndex(ESClient, logger, UserSearchIndex)
	createIndex(ESClient, logger, SeasonSearchIndex)

}

func createIndex(ESClient di.ESClient, logger di.LoggerType, indexType di.ElasticIndex) {
	_, err := esapi.IndicesExistsRequest{
		Index: []string{indexType},
	}.Do(context.Background(), ESClient)

	// If error Index does not exists
	if err != nil {
		ESClient.Indices.Create(indexType)
	}

	logger.Info("Index " + indexType + " Index exits")
}
