package seasons

import (
	"bytes"
	"context"
	"coreservice/internal/di"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type SeasonStatusRepo struct {
	ESClient di.ESClient
	Index    di.ElasticIndex
	Logger   di.LoggerType
}

func New(ESClient di.ESClient, Index di.ElasticIndex, Logger di.LoggerType) *SeasonStatusRepo {
	if ESClient == nil {
		panic("ESClient is nil")
	}

	return &SeasonStatusRepo{
		ESClient: ESClient,
		Index:    Index,
		Logger:   Logger,
	}
}

type SeasonElastic struct {
	SeasonStatus string `json:"season-status"`
}

type SeasonStatus string

const PlannedSeason SeasonStatus = "planned"
const CurrentSeason SeasonStatus = "current"
const EndedSeason SeasonStatus = "ended"

func (sr *SeasonStatusRepo) AddSeasonToIndex(seasonID int) error {

	data, err := json.Marshal(SeasonElastic{string(PlannedSeason)})
	if err != nil {
		sr.Logger.Fatal(err.Error())
		return err
	}

	// making elatic request
	req := esapi.IndexRequest{
		Index:      sr.Index,
		DocumentID: strconv.Itoa(seasonID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// execute it
	resp, err := req.Do(context.Background(), sr.ESClient)
	if err != nil {
		sr.Logger.Fatal(err.Error())
		return err
	}

	defer resp.Body.Close()

	sr.Logger.Info(fmt.Sprintf("Indxed document %s to index %s", resp.String(), sr.Index))
	return nil
}

func (sr *SeasonStatusRepo) StartSeason(seasonID int) error {
	return updateSeasonInIndex(seasonID, CurrentSeason, sr.Logger, sr.Index, sr.ESClient)
}

func (sr *SeasonStatusRepo) EndSeason(seasonID int) error {
	return updateSeasonInIndex(seasonID, EndedSeason, sr.Logger, sr.Index, sr.ESClient)
}

func updateSeasonInIndex(seasonID int, seasonStatus SeasonStatus, logger di.LoggerType, Index di.ElasticIndex, ESClient di.ESClient) error {

	updateData := map[string]interface{}{
		"doc": SeasonElastic{
			SeasonStatus: string(seasonStatus),
		},
	}

	// Преобразуем в JSON
	data, err := json.Marshal(updateData)
	if err != nil {
		logger.Error("Error marshaling update data: " + err.Error())
		return err
	}

	// making elatic request
	req := esapi.UpdateRequest{
		Index:      Index,
		DocumentID: strconv.Itoa(seasonID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// execute it
	resp, err := req.Do(context.Background(), ESClient)
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}

	defer resp.Body.Close()

	logger.Info(fmt.Sprintf("Indxed document %s to index %s", resp.String(), Index))
	return nil
}

func (sr *SeasonStatusRepo) ActiveSeason() ([]int32, error) {
	return searchSeasonsByStatus(CurrentSeason, sr.Logger, sr.Index, sr.ESClient)
}

func (sr *SeasonStatusRepo) PlannedSeasons() ([]int32, error) {
	return searchSeasonsByStatus(PlannedSeason, sr.Logger, sr.Index, sr.ESClient)
}

func searchSeasonsByStatus(status SeasonStatus, logger di.LoggerType, Index di.ElasticIndex, ESClient di.ESClient) ([]int32, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"season-status.keyword": string(status),
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Error("Failed to encode search query: " + err.Error())
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	// Execute query
	res, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(Index), // Используем правильный индекс
		ESClient.Search.WithBody(&buf),
	)
	if err != nil {
		logger.Error("Search request failed: " + err.Error())
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("Elasticsearch error: " + res.String())
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Error(err.Error())
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
