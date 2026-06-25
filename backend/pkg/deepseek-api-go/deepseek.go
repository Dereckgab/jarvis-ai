package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{APIKey: apiKey, BaseURL: baseURL}
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (c *Client) CreateChatCompletion(req ChatRequest) (*ChatResponse, error) {
	data, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(data))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	var chatResp ChatResponse
	json.NewDecoder(resp.Body).Decode(&chatResp)
	return &chatResp, nil
}
