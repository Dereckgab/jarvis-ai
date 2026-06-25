package service

import (
	"context"
)

// AIService defines the interface for AI chat completion services.
type AIService interface {
	ChatCompletion(ctx context.Context, userID string, prompt string) (string, error)
}

// NoOpAIService implements a fallback AI service that returns a placeholder response.
type NoOpAIService struct{}

// NewNoOpAIService creates a new NoOpAIService.
func NewNoOpAIService() AIService {
	return &NoOpAIService{}
}

// ChatCompletion returns a placeholder response when no AI provider is configured.
func (s *NoOpAIService) ChatCompletion(ctx context.Context, userID string, prompt string) (string, error) {
	return "JARVIS AI is not configured. Please set AI_PROVIDER (openai or deepseek) and the corresponding API key in your .env file.", nil
}
