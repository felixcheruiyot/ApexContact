package config

import (
	"fmt"
	"os"
	"strconv"

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
	MediaServerURL  string
	MediaServerKey  string
	RTMPIngestURL   string // public RTMP address broadcasters connect to, e.g. rtmp://example.com/live

	// LiveKit
	LiveKitURL       string // internal Docker URL used by backend → LiveKit server
	LiveKitPublicURL string // public WSS URL sent to browser clients
	LiveKitAPIKey    string
	LiveKitAPISecret string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string

	// App
	AppEnv      string
	AppURL      string
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

		IntaSendPublicKey:   getEnv("INTASEND_PUBLIC_KEY", ""),
		IntaSendPrivateKey:  getEnv("INTASEND_PRIVATE_KEY", ""),
		IntaSendBaseURL:     getEnv("INTASEND_BASE_URL", "https://sandbox.intasend.com"),
		MediaServerURL:      getEnv("MEDIA_SERVER_URL", "http://media-server:8888"),
		MediaServerKey:      getEnv("MEDIA_SERVER_KEY", ""),
		RTMPIngestURL:       getEnv("RTMP_INGEST_URL", "rtmp://localhost/live"),
		LiveKitURL:          getEnv("LIVEKIT_URL", "wss://livekit:7880"),
		LiveKitPublicURL:    getEnv("LIVEKIT_PUBLIC_URL", ""),
		LiveKitAPIKey:       getEnv("LIVEKIT_API_KEY", "devkey"),
		LiveKitAPISecret:    getEnv("LIVEKIT_API_SECRET", "devsecret0000000000000000000000"),
		PostmarkServerToken: getEnv("POSTMARK_SERVER_TOKEN", ""),
		EmailFrom:           getEnv("EMAIL_FROM", ""),
		EmailFromName:       getEnv("EMAIL_FROM_NAME", "Live Streamify"),
		AppEnv:              getEnv("APP_ENV", "development"),
		AppURL:              getEnv("APP_URL", "http://localhost:3000"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
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

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
