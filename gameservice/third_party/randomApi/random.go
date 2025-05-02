package randomapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey string
	apiUrl string
	client *http.Client
}

func NewClient(apiKey string, apiUrl string) *Client {
	return &Client{
		apiKey: apiKey,
		apiUrl: apiUrl,
		client: &http.Client{},
	}
}

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

type generateIntegersResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Random struct {
			Data []int `json:"data"`
		} `json:"random"`
	} `json:"result"`
	ID int `json:"id"`
}

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

	resp, err := c.client.Post(c.apiUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var result generateIntegersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result.Result.Random.Data, nil
}
