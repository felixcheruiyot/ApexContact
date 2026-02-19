package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	// Database
	DatabaseURL string

	// Redis
	RedisAddr     string
	RedisPassword string

	// JWT
	JWTSecret string

	// IntaSend
	IntaSendPublicKey  string
	IntaSendPrivateKey string
	IntaSendBaseURL    string

	// Media server
	MediaServerURL string
	MediaServerKey string

	// LiveKit
	LiveKitURL       string
	LiveKitAPIKey    string
	LiveKitAPISecret string

	// App
	AppEnv      string
	FrontendURL string
}

func Load() (*Config, error) {
	// Load .env in development (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnv("PORT", "8000"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		IntaSendPublicKey:  getEnv("INTASEND_PUBLIC_KEY", ""),
		IntaSendPrivateKey: getEnv("INTASEND_PRIVATE_KEY", ""),
		IntaSendBaseURL:    getEnv("INTASEND_BASE_URL", "https://sandbox.intasend.com"),
		MediaServerURL:     getEnv("MEDIA_SERVER_URL", "http://media-server:8888"),
		MediaServerKey:     getEnv("MEDIA_SERVER_KEY", ""),
		LiveKitURL:         getEnv("LIVEKIT_URL", "ws://livekit:7880"),
		LiveKitAPIKey:      getEnv("LIVEKIT_API_KEY", "devkey"),
		LiveKitAPISecret:   getEnv("LIVEKIT_API_SECRET", "devsecret0000000000000000000000"),
		AppEnv:             getEnv("APP_ENV", "development"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
