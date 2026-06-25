package ai

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/internal/application/service"
	deepseek "jarvis/pkg/deepseek-api-go"
)

// DeepSeekAIService implements the AIService interface for DeepSeek API.
type DeepSeekAIService struct {
	client *deepseek.Client
	cfg    *config.AIConfig
}

// NewDeepSeekAIService creates a new DeepSeekAIService.
func NewDeepSeekAIService(cfg *config.AIConfig) (service.AIService, error) {
	if cfg.DeepSeekAPIKey == "" {
		return nil, fmt.Errorf("deepseek API key is not configured")
	}

	baseURL := cfg.DeepSeekBaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}

	client := deepseek.NewClient(cfg.DeepSeekAPIKey, baseURL)

	return &DeepSeekAIService{
		client: client,
		cfg:    cfg,
	}, nil
}

// ChatCompletion sends a chat completion request to DeepSeek API.
func (s *DeepSeekAIService) ChatCompletion(ctx context.Context, userID string, prompt string) (string, error) {
	resp, err := s.client.CreateChatCompletion(deepseek.ChatRequest{
		Model: s.cfg.Model,
		Messages: []deepseek.Message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get chat completion from DeepSeek: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from DeepSeek API")
	}

	return resp.Choices[0].Message.Content, nil
}
