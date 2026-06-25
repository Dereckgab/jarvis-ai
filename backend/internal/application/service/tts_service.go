package service

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/pkg/tts"
)

// TTSService defines the application service interface for Text-to-Speech operations.
type TTSService interface {
	GenerateSpeech(ctx context.Context, text string) ([]byte, error)
}

// ttsService implements TTSService.
type ttsService struct {
	ttsProvider tts.TTSProvider
	cfg         *config.AIConfig
}

// NewTTSService creates a new TTSService.
func NewTTSService(cfg *config.AIConfig) (TTSService, error) {
	provider, err := tts.NewTTSProvider(cfg.TTSProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TTS provider: %w", err)
	}
	return &ttsService{
		ttsProvider: provider,
		cfg:         cfg,
	}, nil
}

// GenerateSpeech generates speech from the given text using the configured TTS provider.
func (s *ttsService) GenerateSpeech(ctx context.Context, text string) ([]byte, error) {
	// TODO: Add caching logic here using Redis if needed
	return s.ttsProvider.GenerateSpeech(ctx, text)
}
