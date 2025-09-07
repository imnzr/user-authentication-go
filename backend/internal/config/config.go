package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       ServerConfig   `json:"server"`
	Database     DatabaseConfig `json:"database"`
	Logger       LoggerConfig   `json:"logger"`
	JSONWebToken JWTConfig      `json:"json_web_token"`
	RedisCfg     RedisConfig
}

type ServerConfig struct {
	Port         int           `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

type LoggerConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

type JWTConfig struct {
	JWTSecretKey         string        `json:"jwt_secret_key"`
	AccessTokenDuration  time.Duration `json:"access_token"`
	RefreshTokenDuration time.Duration `json:"refresh_token"`
}

type RedisConfig struct {
	DBUrl     string
	RedisAddr string
	RedisPass string
	RedisDB   int
}

// Load configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env doesn't exist, just continue
		fmt.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{}

	// Load server config
	cfg.Server = ServerConfig{
		Port:         getEnvIntOrDefault("SERVER_PORT", 8080),
		Host:         getEnvOrDefault("SERVER_HOST", "localhost"),
		ReadTimeout:  getEnvDurationOrDefault("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getEnvDurationOrDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  getEnvDurationOrDefault("SERVER_IDLE_TIMEOUT", 60*time.Second),
	}

	// Load database config
	if err := loadDatabaseConfig(&cfg.Database); err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Load logger config
	cfg.Logger = LoggerConfig{
		Level:  getEnvOrDefault("LOG_LEVEL", "info"),
		Format: getEnvOrDefault("LOG_FORMAT", "json"),
	}

	// Load JWT config
	cfg.JSONWebToken = JWTConfig{
		JWTSecretKey:         os.Getenv("JWT_SECRET_KEY"),
		AccessTokenDuration:  getEnvDurationOrDefault("ACCESS_TOKEN", 30*time.Second),
		RefreshTokenDuration: getEnvDurationOrDefault("REFRESH_TOKEN", 60*time.Second),
	}

	// Load Redis Config
	cfg.RedisCfg = RedisConfig{
		DBUrl:     os.Getenv("REDIS_URL"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB:   getEnvIntOrDefault("REDIS_DB", 0),
	}

	return cfg, nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
