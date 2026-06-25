package ai

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/internal/application/service"

	openai "github.com/sashabaranov/go-openai"
)

const groqBaseURL = "https://api.groq.com/openai/v1"

// GroqAIService implements AIService using Groq's OpenAI-compatible API.
type GroqAIService struct {
	client *openai.Client
	cfg    *config.AIConfig
}

// NewGroqAIService creates a new GroqAIService.
func NewGroqAIService(cfg *config.AIConfig) (service.AIService, error) {
	if cfg.GroqAPIKey == "" {
		return nil, fmt.Errorf("groq API key is not configured")
	}

	clientConfig := openai.DefaultConfig(cfg.GroqAPIKey)
	clientConfig.BaseURL = groqBaseURL

	return &GroqAIService{
		client: openai.NewClientWithConfig(clientConfig),
		cfg:    cfg,
	}, nil
}

// ChatCompletion sends a chat completion request to Groq API.
func (s *GroqAIService) ChatCompletion(ctx context.Context, userID string, prompt string) (string, error) {
	model := s.cfg.Model
	if model == "" || model == "deepseek-chat" {
		model = "llama-3.3-70b-versatile"
	}

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are JARVIS, an advanced AI assistant. CRITICAL RULE: Always respond in the exact same language the user writes in. If the user writes in Portuguese, respond entirely in Portuguese. If in English, respond in English. Never switch languages. Be helpful, concise, and intelligent."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: float32(s.cfg.Temperature),
		MaxTokens:   s.cfg.MaxTokens,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get chat completion from Groq: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from Groq API")
	}

	return resp.Choices[0].Message.Content, nil
}
