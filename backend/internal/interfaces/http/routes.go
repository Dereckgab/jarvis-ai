package http

import (
	"jarvis/config"
	"jarvis/internal/interfaces/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"strings"
	"time"
)

// SetupRoutes registers all application routes and middleware.
func SetupRoutes(
	app *fiber.App,
	userHandler *UserHandler,
	aiHandler *AIHandler,
	systemInfoHandler *SystemInfoHandler,
	memoryHandler *MemoryHandler,
	gameHandler *GameHandler,
	ttsHandler *TTSHandler,
	cfg *config.Config,
) {
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Security.RateLimitMax,
		Expiration: cfg.Security.RateLimitWindow,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(ErrorResponse{
				Message: "too many requests",
				Details: "rate limit exceeded, please slow down",
			})
		},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"service":   "jarvis-backend",
			"timestamp": time.Now().UTC(),
		})
	})

	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.RegisterUser)
	auth.Post("/login", userHandler.LoginUser)
	auth.Post("/refresh", userHandler.RefreshToken)

	// All routes below require authentication
	authMiddleware := middleware.AuthRequired(cfg)

	// User profile
	user := api.Group("/user", authMiddleware)
	user.Get("/profile", userHandler.GetUserProfile)
	user.Put("/profile", userHandler.UpdateUserProfile)
	user.Put("/password", userHandler.ChangePassword)

	// AI chat
	ai := api.Group("/ai", authMiddleware)
	ai.Post("/chat", aiHandler.ChatCompletion)

	// System information
	sysInfo := api.Group("/system-info", authMiddleware)
	sysInfo.Post("/collect", systemInfoHandler.CollectSystemInfo)
	sysInfo.Get("/latest", systemInfoHandler.GetLatestSystemInfo)
	sysInfo.Get("/history", systemInfoHandler.GetSystemInfoHistory)

	// Memory
	memory := api.Group("/memory", authMiddleware)
	memory.Post("", memoryHandler.SaveMemory)
	memory.Get("", memoryHandler.GetMemories)
	memory.Post("/search", memoryHandler.SearchMemories)
	memory.Delete("/:id", memoryHandler.DeleteMemory)

	// Games
	games := api.Group("/games", authMiddleware)
	games.Get("/search", gameHandler.SearchGames)
	games.Get("/:id", gameHandler.GetGame)
	games.Post("", gameHandler.CreateGame)
	games.Put("/:id", gameHandler.UpdateGame)
	games.Delete("/:id", gameHandler.DeleteGame)
	games.Get("/:id/requirements", gameHandler.GetGameRequirements)
	games.Post("/:id/requirements", gameHandler.AddGameRequirement)
	games.Put("/:id/requirements/:reqID", gameHandler.UpdateGameRequirement)
	games.Delete("/:id/requirements/:reqID", gameHandler.DeleteGameRequirement)

	// TTS (only if service is available)
	if ttsHandler != nil {
		tts := api.Group("/tts", authMiddleware)
		tts.Post("/generate", ttsHandler.GenerateSpeech)
	}

	// 404 fallback
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Message: "route not found",
			Details: c.Path(),
		})
	})
}
