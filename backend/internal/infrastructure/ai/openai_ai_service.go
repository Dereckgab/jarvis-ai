package ai

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/internal/application/service"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIService implements the AIService interface for OpenAI API.
type OpenAIService struct {
	client *openai.Client
	cfg    *config.AIConfig
}

// NewOpenAIService creates a new OpenAIService.
func NewOpenAIService(cfg *config.AIConfig) (service.AIService, error) {
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("openai API key is not configured")
	}

	clientConfig := openai.DefaultConfig(cfg.OpenAIAPIKey)
	if cfg.OpenAIBaseURL != "" {
		clientConfig.BaseURL = cfg.OpenAIBaseURL
	}
	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIService{
		client: client,
		cfg:    cfg,
	}, nil
}

// ChatCompletion sends a chat completion request to OpenAI API.
func (s *OpenAIService) ChatCompletion(ctx context.Context, userID string, prompt string) (string, error) {
	// For now, userID is not used by OpenAI API directly, but can be used for logging or context management.
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.cfg.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: float32(s.cfg.Temperature),
		MaxTokens:   s.cfg.MaxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get chat completion from OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI API")
	}

	return resp.Choices[0].Message.Content,
		nil
}
