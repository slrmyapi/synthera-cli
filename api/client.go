package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

)

var baseURL = os.Getenv("BASE_URL")

func init() {
	if baseURL == "" {
		baseURL = "https://workers.slrmyapi.workers.dev"
	}
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiToken: token,
	}
}

func (c *Client) makeRequest(method, endpoint string, requestBody, response any) error {
	var reqBody io.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, endpoint), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute requst: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-OK status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}


func (c *Client) TraceName(name string) ([]TraceNameItem, error) {
	req := TraceNameRequest{
		Name: name,
	}
	var res TraceNameResponse
	err := c.makeRequest("POST", "/trace/name", req, &res)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) TraceDetail(id int) ([]TraceDetailID, error) {
	req := TraceDetailRequest{
		ID: id,
	}
	var res TraceDetailResponse
	err := c.makeRequest("POST", "/trace/id", req, &res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) TraceRelations(id int, offset int) ([]TraceDetailID, TraceRelationsItem, error) {
	req := TraceRelationsRequest{
		ID: id,
		Offset: offset,
	}

	var res TraceRelationsResponse
	err := c.makeRequest("POST", "/trace/relations", req, &res)
	if err != nil {
		return nil, TraceRelationsItem{}, err
	}

	return res.Data, res.Relationships, nil
}
