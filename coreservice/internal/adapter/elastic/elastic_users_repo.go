package elastic

// import (
// 	"bytes"
// 	"context"
// 	"coreservice/internal/entity"
// 	"coreservice/internal/logger"
// 	"coreservice/internal/repository/sqlc"
// 	"coreservice/pkg"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/elastic/go-elasticsearch/v8/esapi"
// )

// type SearchType string

// const Strict SearchType = "strict"
// const Wildcard SearchType = "wildcard"
// const Fuzzy SearchType = "fuzzy"

// func AddingAllUsersToIndex() ([]entity.User, error) {
// 	// Recieving all users form DB
// 	users, err := sqlc.GetAllUsers()
// 	if err != nil {
// 		logger.Logger.Fatal(err.Error())
// 		return []entity.User{}, err
// 	}

// 	response := make([]entity.User, 0, len(users))

// 	for _, user := range users {
// 		logger.Logger.Info(fmt.Sprint(user))

// 		// Convert user from DB to strcut format
// 		res := pkg.GetUserInfo(&user)
// 		response = append(response, res)

// 		if err := AddUserToIndex(res, int(user.ID)); err != nil {
// 			return []entity.User{}, err
// 		}
// 	}
// 	return response, nil
// }

// func AddUserToIndex(user entity.User, userId int) error {

// 	data, err := json.Marshal(user)
// 	if err != nil {
// 		logger.Logger.Fatal(err.Error())
// 		return err
// 	}

// 	// making elatic request
// 	req := esapi.IndexRequest{
// 		Index:      UserSearchIndex,
// 		DocumentID: strconv.Itoa(userId),
// 		Body:       bytes.NewReader(data),
// 		Refresh:    "true",
// 	}

// 	// execute it
// 	resp, err := req.Do(context.Background(), ESClient)
// 	if err != nil {
// 		logger.Logger.Fatal(err.Error())
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	logger.Logger.Info(fmt.Sprintf("Indxed document %s to index %s", resp.String(), UserSearchIndex))
// 	return nil
// }

// func GetUserByName(name string, searchType SearchType) ([]entity.User, error) {
// 	if name == "" {
// 		return []entity.User{}, nil
// 	}

// 	query, err := createUserQuery(name, searchType)
// 	if err != nil {
// 		return []entity.User{}, nil
// 	}

// 	// Marshalling query into json
// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(query); err != nil {
// 		logger.Logger.Error(err.Error())
// 		return []entity.User{}, err
// 	}

// 	// execute query
// 	res, err := ESClient.Search(
// 		ESClient.Search.WithIndex(UserSearchIndex),
// 		ESClient.Search.WithBody(&buf),
// 	)
// 	defer res.Body.Close()

// 	if err != nil || res.IsError() {
// 		logger.Logger.Fatal(err.Error())
// 		return []entity.User{}, err
// 	}

// 	// decoding responce from elastic
// 	var r map[string]interface{}
// 	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 		logger.Logger.Error(err.Error())
// 		return []entity.User{}, err
// 	}

// 	// Extracting id's from response to find them in DB
// 	var ids []int32
// 	if hits, ok := r["hits"].(map[string]interface{}); ok {
// 		if hitsHits, ok := hits["hits"].([]interface{}); ok {
// 			for _, hit := range hitsHits {
// 				if hitMap, ok := hit.(map[string]interface{}); ok {
// 					if idStr, ok := hitMap["_id"].(string); ok {
// 						id, _ := strconv.Atoi(idStr)
// 						ids = append(ids, int32(id))
// 					}
// 				}
// 			}
// 		}
// 	}

// 	users, err := sqlc.GetUsersByIds(ids)
// 	if err != nil {
// 		logger.Logger.Error(err.Error())
// 		return []entity.User{}, err
// 	}

// 	return pkg.ConvertDBUserSliceToUser(users), nil

// }

// func createUserQuery(name string, searchType SearchType) (map[string]interface{}, error) {
// 	switch searchType {
// 	case Strict:
// 		return createStrictQuery(name), nil
// 	case Wildcard:
// 		return createWildcardQuery(name), nil
// 	case Fuzzy:
// 		return createFuzzyQuery(name), nil
// 	default:
// 		return nil, errors.New("there is no such search type")
// 	}
// }

// func createStrictQuery(name string) map[string]interface{} {
// 	return map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"multi_match": map[string]interface{}{
// 				"query":  name,
// 				"fields": []string{"login"},
// 			},
// 		},
// 	}
// }

// func createWildcardQuery(name string) map[string]interface{} {
// 	return map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"wildcard": map[string]interface{}{
// 				"login": map[string]interface{}{
// 					"value":            "*" + strings.ToLower(name) + "*", // ищет подстроку (регистронезависимо)
// 					"case_insensitive": true,
// 				},
// 			},
// 		},
// 	}
// }

// func createFuzzyQuery(name string) map[string]interface{} {
// 	return map[string]interface{}{
// 		"query": map[string]interface{}{
// 			"fuzzy": map[string]interface{}{
// 				"login": map[string]interface{}{
// 					"value":     name,
// 					"fuzziness": 2,
// 					// "fuzziness":"AUTO",
// 				},
// 			},
// 		},
// 	}
// }
