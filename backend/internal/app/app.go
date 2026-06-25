package app

import (
	"context"
	"fmt"

	"jarvis/config"
	"jarvis/internal/application/service"
	infraAI "jarvis/internal/infrastructure/ai"
	"jarvis/internal/infrastructure/database"
	infraRepo "jarvis/internal/infrastructure/repository"
	httpHandlers "jarvis/internal/interfaces/http"
	"jarvis/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	Fiber *fiber.App
	Cfg   *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	// Database
	db, err := database.NewMySQLConnection(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	logger.Info("Database connected")

	if err := database.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Database migrations applied")

	// AI service
	var aiSvc service.AIService
	switch cfg.AI.Provider {
	case "openai":
		aiSvc, err = infraAI.NewOpenAIService(&cfg.AI)
		if err != nil {
			logger.Warn("Failed to initialize OpenAI service, falling back to NoOp", "error", err)
			aiSvc = service.NewNoOpAIService()
		} else {
			logger.Info("AI provider: OpenAI", "model", cfg.AI.Model)
		}
	case "deepseek":
		aiSvc, err = infraAI.NewDeepSeekAIService(&cfg.AI)
		if err != nil {
			logger.Warn("Failed to initialize DeepSeek service, falling back to NoOp", "error", err)
			aiSvc = service.NewNoOpAIService()
		} else {
			logger.Info("AI provider: DeepSeek", "model", cfg.AI.Model)
		}
	case "groq":
		aiSvc, err = infraAI.NewGroqAIService(&cfg.AI)
		if err != nil {
			logger.Warn("Failed to initialize Groq service, falling back to NoOp", "error", err)
			aiSvc = service.NewNoOpAIService()
		} else {
			logger.Info("AI provider: Groq", "model", cfg.AI.Model)
		}
	default:
		logger.Warn("No AI provider configured, using NoOp", "provider", cfg.AI.Provider)
		aiSvc = service.NewNoOpAIService()
	}

	// TTS service (non-fatal if unavailable)
	ttsSvc, err := service.NewTTSService(&cfg.AI)
	if err != nil {
		logger.Warn("TTS service unavailable", "error", err)
		ttsSvc = nil
	} else {
		logger.Info("TTS provider initialized", "provider", cfg.AI.TTSProvider)
	}

	// Repositories
	userRepo := infraRepo.NewGormUserRepository(db)
	gameRepo := infraRepo.NewGormGameRepository(db)
	sysInfoRepo := infraRepo.NewGormSystemInfoRepository(db)
	memoryRepo := infraRepo.NewInMemoryMemoryRepository()

	// Services
	userSvc := service.NewUserService(userRepo)
	gameSvc := service.NewGameService(gameRepo)
	sysInfoSvc := service.NewSystemInfoService(sysInfoRepo)
	memorySvc := service.NewMemoryService(memoryRepo, aiSvc)

	// Handlers
	userHandler := httpHandlers.NewUserHandler(userSvc, cfg)
	aiHandler := httpHandlers.NewAIHandler(aiSvc, sysInfoSvc, cfg)
	sysInfoHandler := httpHandlers.NewSystemInfoHandler(sysInfoSvc, cfg)
	memoryHandler := httpHandlers.NewMemoryHandler(memorySvc, cfg)
	gameHandler := httpHandlers.NewGameHandler(gameSvc, cfg)

	var ttsHandler *httpHandlers.TTSHandler
	if ttsSvc != nil {
		ttsHandler = httpHandlers.NewTTSHandler(ttsSvc, cfg)
	}

	// Fiber app
	f := fiber.New(fiber.Config{
		ErrorHandler: httpHandlers.ErrorHandler,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
	})

	httpHandlers.SetupRoutes(f, userHandler, aiHandler, sysInfoHandler, memoryHandler, gameHandler, ttsHandler, cfg)

	return &App{
		Fiber: f,
		Cfg:   cfg,
	}, nil
}

// Run starts the HTTP server. Blocks until the server stops.
func (a *App) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", a.Cfg.App.Port) // fixed: was %s but Port is int
	logger.Info("JARVIS server starting", "addr", addr, "env", a.Cfg.App.Env)
	return a.Fiber.Listen(addr)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.Fiber.ShutdownWithContext(ctx)
}
