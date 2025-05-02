package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	apiKey := "ff9512b9-99c4-4b88-bf63-a40b27e13bb8"
	client := NewClient(apiKey)

	// Генерация 10 чисел от 1 до 100
	numbers, err := client.GenerateIntegers(2, 1, 6, true)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Printf("Сгенерированные числа: %v\n", numbers)
}

const apiURL = "https://api.random.org/json-rpc/2/invoke"

type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Request структура для запроса
type generateIntegersRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		APIKey      string `json:"apiKey"`
		N           int    `json:"n"`
		Min         int    `json:"min"`
		Max         int    `json:"max"`
		Replacement bool   `json:"replacement"`
	} `json:"params"`
	ID int `json:"id"`
}

// Response структура для ответа
type generateIntegersResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Random struct {
			Data []int `json:"data"`
		} `json:"random"`
	} `json:"result"`
	ID int `json:"id"`
}

// GenerateIntegers генерирует случайные целые числа
func (c *Client) GenerateIntegers(n, min, max int, replacement bool) ([]int, error) {
	reqBody := generateIntegersRequest{
		JSONRPC: "2.0",
		Method:  "generateIntegers",
		ID:      1,
	}
	reqBody.Params.APIKey = c.apiKey
	reqBody.Params.N = n
	reqBody.Params.Min = min
	reqBody.Params.Max = max
	reqBody.Params.Replacement = replacement

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.client.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var result generateIntegersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result.Result.Random.Data, nil
}
