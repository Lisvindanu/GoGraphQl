package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	Environment  string
	RateLimitRPM int
	CORSOrigins  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Environment: getEnv("ENVIRONMENT", "dev"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),
	}

	rpm, err := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "60"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPM: %w", err)
	}
	cfg.RateLimitRPM = rpm

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func (c *Config) IsDev() bool {
	return c.Environment == "dev"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
