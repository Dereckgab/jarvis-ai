package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"jarvis/config"
	"jarvis/internal/app"
	"jarvis/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.InitLogger("development")
		logger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.App.Env)

	// Create new application instance
	application, err := app.NewApp(cfg)
	if err != nil {
		logger.Fatal("Failed to create application", "error", err)
	}

	// Start the application in a goroutine
	go func() {
		logger.Info("Starting application server...")
		if err := application.Run(context.Background()); err != nil {
			logger.Error("Application server failed", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		logger.Fatal("Application shutdown failed", "error", err)
	}

	logger.Info("Application gracefully stopped.")
}
