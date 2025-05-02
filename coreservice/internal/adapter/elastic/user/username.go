package user

import (
	"bytes"
	"context"
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	userRepo "coreservice/internal/adapter/db/postgres/user"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type UserNameRepo struct {
	ESClient di.ESClient
	Index    entity.ElasticIndexType
	Logger   di.LoggerType
	Repo     *userRepo.UserRepo
}

func New(ESClient di.ESClient, Index entity.ElasticIndexType, Logger di.LoggerType, Repo *userRepo.UserRepo) *UserNameRepo {
	if ESClient == nil {
		panic("ESClient is nil")
	}

	return &UserNameRepo{
		ESClient: ESClient,
		Index:    Index,
		Logger:   Logger,
		Repo:     Repo,
	}
}

type SearchType string

const Strict SearchType = "strict"
const Wildcard SearchType = "wildcard"
const Fuzzy SearchType = "fuzzy"

func (er *UserNameRepo) AddingAllUsersToIndex() ([]entity.User, error) {
	// Recieving all users form DB
	users, err := er.Repo.GetAllUsers()
	if err != nil {
		er.Logger.Fatal(err.Error())
		return []entity.User{}, err
	}

	response := make([]entity.User, 0, len(users))

	for _, user := range users {
		er.Logger.Info(fmt.Sprint(user))

		// Convert user from DB to strcut format
		res := pkg.GetUserInfo(&user)
		response = append(response, res)

		if err := er.AddUserToIndex(res, int(user.ID)); err != nil {
			return []entity.User{}, err
		}
	}
	return response, nil
}

func (er *UserNameRepo) AddUserToIndex(user entity.User, userId int) error {

	data, err := json.Marshal(user)
	if err != nil {
		er.Logger.Fatal(err.Error())
		return err
	}

	// making elatic request
	req := esapi.IndexRequest{
		Index:      string(er.Index),
		DocumentID: strconv.Itoa(userId),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// execute it
	resp, err := req.Do(context.Background(), er.ESClient)
	if err != nil {
		er.Logger.Fatal(err.Error())
		return err
	}

	defer resp.Body.Close()

	er.Logger.Info(fmt.Sprintf("Indxed document %s to index %s", resp.String(), er.Index))
	return nil
}

func (er *UserNameRepo) GetUserByNameStrict(name string) ([]entity.User, error) {
	return getUserByName(name, Strict, er.Logger, er.ESClient, er.Index, er.Repo)
}

func (er *UserNameRepo) GetUserByNameWildcard(name string) ([]entity.User, error) {
	return getUserByName(name, Wildcard, er.Logger, er.ESClient, er.Index, er.Repo)
}

func (er *UserNameRepo) GetUserByNameFuzzy(name string) ([]entity.User, error) {
	return getUserByName(name, Fuzzy, er.Logger, er.ESClient, er.Index, er.Repo)
}

func getUserByName(name string, searchType SearchType, logger di.LoggerType, ESClient di.ESClient, Index entity.ElasticIndexType, repo *userRepo.UserRepo) ([]entity.User, error) {
	if name == "" {
		return []entity.User{}, nil
	}

	query, err := createUserQuery(name, searchType)
	if err != nil {
		return []entity.User{}, nil
	}

	// Marshalling query into json
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Error(err.Error())
		return []entity.User{}, err
	}

	// execute query
	res, err := ESClient.Search(
		ESClient.Search.WithIndex(string(Index)),
		ESClient.Search.WithBody(&buf),
	)
	defer res.Body.Close()

	if err != nil || res.IsError() {
		logger.Fatal(err.Error())
		return []entity.User{}, err
	}

	// decoding responce from elastic
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Error(err.Error())
		return []entity.User{}, err
	}

	// Extracting id's from response to find them in DB
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

	users, err := repo.GetUsersByIds(ids)
	if err != nil {
		logger.Error(err.Error())
		return []entity.User{}, err
	}

	return pkg.ConvertDBUserSliceToUser(users), nil

}

func createUserQuery(name string, searchType SearchType) (map[string]interface{}, error) {
	switch searchType {
	case Strict:
		return createStrictQuery(name), nil
	case Wildcard:
		return createWildcardQuery(name), nil
	case Fuzzy:
		return createFuzzyQuery(name), nil
	default:
		return nil, errors.New("there is no such search type")
	}
}

func createStrictQuery(name string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  name,
				"fields": []string{"login"},
			},
		},
	}
}

func createWildcardQuery(name string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"login": map[string]interface{}{
					"value":            "*" + strings.ToLower(name) + "*", // ищет подстроку (регистронезависимо)
					"case_insensitive": true,
				},
			},
		},
	}
}

func createFuzzyQuery(name string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy": map[string]interface{}{
				"login": map[string]interface{}{
					"value":     name,
					"fuzziness": 2,
					// "fuzziness":"AUTO",
				},
			},
		},
	}
}
