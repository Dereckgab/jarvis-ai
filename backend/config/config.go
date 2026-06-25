package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:"app" validate:"required"`
	Database  DatabaseConfig  `mapstructure:"database" validate:"required"`
	Redis     RedisConfig     `mapstructure:"redis" validate:"required"`
	Security  SecurityConfig  `mapstructure:"security" validate:"required"`
	AI        AIConfig        `mapstructure:"ai" validate:"required"`
	Qdrant    QdrantConfig    `mapstructure:"qdrant" validate:"required"`
	Telemetry TelemetryConfig `mapstructure:"telemetry" validate:"required"`
}

type AppConfig struct {
	Env             string        `mapstructure:"env" validate:"required,oneof=development staging production"`
	Port            int           `mapstructure:"port" validate:"required,gt=0"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" validate:"required"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required,gt=0"`
	User            string        `mapstructure:"user" validate:"required"`
	Password        string        `mapstructure:"password" validate:"required"`
	Name            string        `mapstructure:"name" validate:"required"`
	Charset         string        `mapstructure:"charset" validate:"required"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"required,gt=0"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"required,gt=0"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"required"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db" validate:"gte=0"`
}

type SecurityConfig struct {
	JWTSecret           string        `mapstructure:"jwt_secret" validate:"required,min=32"`
	JWTExpiresIn        time.Duration `mapstructure:"jwt_expires_in" validate:"required"`
	JWTRefreshExpiresIn time.Duration `mapstructure:"jwt_refresh_expires_in" validate:"required"`
	RateLimitMax        int           `mapstructure:"rate_limit_max" validate:"required,gt=0"`
	RateLimitWindow     time.Duration `mapstructure:"rate_limit_window" validate:"required"`
	CORSOrigins         []string      `mapstructure:"cors_origins" validate:"required,min=1,dive,uri"`
}

type AIConfig struct {
	Provider          string  `mapstructure:"provider" validate:"required,oneof=deepseek openai groq noop"`
	DeepSeekAPIKey    string  `mapstructure:"deepseek_api_key"`
	DeepSeekBaseURL   string  `mapstructure:"deepseek_base_url"`
	OpenAIAPIKey      string  `mapstructure:"openai_api_key"`
	OpenAIBaseURL     string  `mapstructure:"openai_base_url"`
	GroqAPIKey        string  `mapstructure:"groq_api_key"`
	Model             string  `mapstructure:"model" validate:"required"`
	Temperature       float64 `mapstructure:"temperature" validate:"gte=0,lte=1"`
	MaxTokens         int     `mapstructure:"max_tokens" validate:"gt=0"`
	TTSProvider       string  `mapstructure:"tts_provider" validate:"required,oneof=piper elevenlabs os_native"`
	ElevenLabsAPIKey  string  `mapstructure:"elevenlabs_api_key"`
	ElevenLabsVoiceID string  `mapstructure:"elevenlabs_voice_id"`
}

type QdrantConfig struct {
	Host           string `mapstructure:"host" validate:"required"`
	Port           int    `mapstructure:"port" validate:"required,gt=0"`
	APIKey         string `mapstructure:"api_key"`
	CollectionName string `mapstructure:"collection_name" validate:"required"`
}

type TelemetryConfig struct {
	CollectorEndpoint string `mapstructure:"collector_endpoint" validate:"required"`
	ServiceName       string `mapstructure:"service_name" validate:"required"`
	ServiceVersion    string `mapstructure:"service_version" validate:"required"`
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists in the backend directory or repo root
	_ = godotenv.Load(".env", "../.env", "./backend/.env")

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // Read environment variables

	// Bind environment variables explicitly to nested configuration keys.
	// This ensures Viper correctly resolves keys like redis.password.
	_ = v.BindEnv("app.env", "APP_ENV")
	_ = v.BindEnv("app.port", "APP_PORT")
	_ = v.BindEnv("app.read_timeout", "APP_READ_TIMEOUT")
	_ = v.BindEnv("app.write_timeout", "APP_WRITE_TIMEOUT")
	_ = v.BindEnv("app.shutdown_timeout", "APP_SHUTDOWN_TIMEOUT")

	_ = v.BindEnv("database.host", "DATABASE_HOST")
	_ = v.BindEnv("database.port", "DATABASE_PORT")
	_ = v.BindEnv("database.user", "DATABASE_USER")
	_ = v.BindEnv("database.password", "DATABASE_PASSWORD")
	_ = v.BindEnv("database.name", "DATABASE_NAME")
	_ = v.BindEnv("database.charset", "DATABASE_CHARSET")
	_ = v.BindEnv("database.max_idle_conns", "DATABASE_MAX_IDLE_CONNS")
	_ = v.BindEnv("database.max_open_conns", "DATABASE_MAX_OPEN_CONNS")
	_ = v.BindEnv("database.conn_max_lifetime", "DATABASE_CONN_MAX_LIFETIME")

	_ = v.BindEnv("redis.addr", "REDIS_ADDR")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = v.BindEnv("redis.db", "REDIS_DB")

	_ = v.BindEnv("security.jwt_secret", "SECURITY_JWT_SECRET")
	_ = v.BindEnv("security.jwt_expires_in", "SECURITY_JWT_EXPIRES_IN")
	_ = v.BindEnv("security.jwt_refresh_expires_in", "SECURITY_JWT_REFRESH_EXPIRES_IN")
	_ = v.BindEnv("security.rate_limit_max", "SECURITY_RATE_LIMIT_MAX")
	_ = v.BindEnv("security.rate_limit_window", "SECURITY_RATE_LIMIT_WINDOW")
	_ = v.BindEnv("security.cors_origins", "SECURITY_CORS_ORIGINS")

	_ = v.BindEnv("ai.provider", "AI_PROVIDER")
	_ = v.BindEnv("ai.deepseek_api_key", "AI_DEEPSEEK_API_KEY")
	_ = v.BindEnv("ai.deepseek_base_url", "AI_DEEPSEEK_BASE_URL")
	_ = v.BindEnv("ai.openai_api_key", "AI_OPENAI_API_KEY")
	_ = v.BindEnv("ai.openai_base_url", "AI_OPENAI_BASE_URL")
	_ = v.BindEnv("ai.groq_api_key", "AI_GROQ_API_KEY")
	_ = v.BindEnv("ai.model", "AI_MODEL")
	_ = v.BindEnv("ai.temperature", "AI_TEMPERATURE")
	_ = v.BindEnv("ai.max_tokens", "AI_MAX_TOKENS")
	_ = v.BindEnv("ai.tts_provider", "AI_TTS_PROVIDER")
	_ = v.BindEnv("ai.elevenlabs_api_key", "AI_ELEVENLABS_API_KEY")
	_ = v.BindEnv("ai.elevenlabs_voice_id", "AI_ELEVENLABS_VOICE_ID")

	_ = v.BindEnv("qdrant.host", "QDRANT_HOST")
	_ = v.BindEnv("qdrant.port", "QDRANT_PORT")
	_ = v.BindEnv("qdrant.api_key", "QDRANT_API_KEY")
	_ = v.BindEnv("qdrant.collection_name", "QDRANT_COLLECTION_NAME")

	_ = v.BindEnv("telemetry.collector_endpoint", "TELEMETRY_COLLECTOR_ENDPOINT")
	_ = v.BindEnv("telemetry.service_name", "TELEMETRY_SERVICE_NAME")
	_ = v.BindEnv("telemetry.service_version", "TELEMETRY_SERVICE_VERSION")

	// Set default values
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.read_timeout", "5s")
	v.SetDefault("app.write_timeout", "10s")
	v.SetDefault("app.shutdown_timeout", "15s")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.name", "jarvis_fullia")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.db", 0)

	v.SetDefault("security.jwt_secret", "supersecretjwtkeythatshouldbemorethan32chars")
	v.SetDefault("security.jwt_expires_in", "15m")
	v.SetDefault("security.jwt_refresh_expires_in", "72h")
	v.SetDefault("security.rate_limit_max", 100)
	v.SetDefault("security.rate_limit_window", "1m")
	v.SetDefault("security.cors_origins", []string{"http://localhost:3000", "http://localhost:8080"})

	v.SetDefault("ai.provider", "noop")
	v.SetDefault("ai.model", "deepseek-chat")
	v.SetDefault("ai.temperature", 0.7)
	v.SetDefault("ai.max_tokens", 1024)
	v.SetDefault("ai.tts_provider", "piper")

	v.SetDefault("qdrant.host", "localhost")
	v.SetDefault("qdrant.port", 6333)
	v.SetDefault("qdrant.collection_name", "jarvis_memory")

	v.SetDefault("telemetry.collector_endpoint", "localhost:4317")
	v.SetDefault("telemetry.service_name", "jarvis-backend")
	v.SetDefault("telemetry.service_version", "1.0.0")

	if origins := v.GetString("security.cors_origins"); origins != "" {
		v.Set("security.cors_origins", strings.Split(origins, ","))
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// GetDSN returns the DSN for the database connection
func (dc *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dc.User, dc.Password, dc.Host, dc.Port, dc.Name, dc.Charset)
}
