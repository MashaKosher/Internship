package elastic

import (
	"bytes"
	"context"
	"coreservice/internal/logger"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type SeasonElastic struct {
	SeasonStatus string `json:"season-status"`
}

type SeasonStatus string

const PlannedSeason SeasonStatus = "planned"
const CurrentSeason SeasonStatus = "current"
const EndedSeason SeasonStatus = "ended"

func AddSeasonToIndex(seasonID int) error {

	data, err := json.Marshal(SeasonElastic{string(PlannedSeason)})
	if err != nil {
		logger.Logger.Fatal(err.Error())
		return err
	}

	// making elatic request
	req := esapi.IndexRequest{
		Index:      SeasonSearchIndex,
		DocumentID: strconv.Itoa(seasonID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// execute it
	resp, err := req.Do(context.Background(), ESClient)
	if err != nil {
		logger.Logger.Fatal(err.Error())
		return err
	}

	defer resp.Body.Close()

	logger.Logger.Info(fmt.Sprintf("Indxed document %s to index %s", resp.String(), UserSearchIndex))
	return nil
}

// func UpdateSeasonInIndex(sea)

func UpdateSeasonInIndex(seasonID int, seasonStatus SeasonStatus) error {

	updateData := map[string]interface{}{
		"doc": SeasonElastic{
			SeasonStatus: string(seasonStatus),
		},
	}

	// Преобразуем в JSON
	data, err := json.Marshal(updateData)
	if err != nil {
		logger.Logger.Error("Error marshaling update data: " + err.Error())
		return err
	}

	// making elatic request
	req := esapi.UpdateRequest{
		Index:      SeasonSearchIndex,
		DocumentID: strconv.Itoa(seasonID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// execute it
	resp, err := req.Do(context.Background(), ESClient)
	if err != nil {
		logger.Logger.Fatal(err.Error())
		return err
	}

	defer resp.Body.Close()

	logger.Logger.Info(fmt.Sprintf("Indxed document %s to index %s", resp.String(), UserSearchIndex))
	return nil
}

// SearchSeasonsByStatus выполняет поиск сезонов по статусу
// SearchSeasonsByStatus выполняет поиск сезонов по статусу
func SearchSeasonsByStatus(status SeasonStatus) ([]int32, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"season-status.keyword": string(status),
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Logger.Error("Failed to encode search query: " + err.Error())
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	// Execute query
	res, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(SeasonSearchIndex), // Используем правильный индекс
		ESClient.Search.WithBody(&buf),
	)
	if err != nil {
		logger.Logger.Error("Search request failed: " + err.Error())
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Logger.Error("Elasticsearch error: " + res.String())
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	var ids []int32
	if hits, ok := r["hits"].(map[string]interface{}); ok {
		if hitsHits, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsHits {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if idStr, ok := hitMap["_id"].(string); ok {
						id, _ := strconv.Atoi(idStr)
						ids = append(ids, int32(id))
					}
				}
			}
		}
	}

	return ids, nil
}
