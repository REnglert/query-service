package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceName string
	Port        string
	LLMBaseURL  string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		ServiceName: "query-service",
		Port:        port,
		LLMBaseURL:  getEnv("LLM_BASE_URL", "http://localhost:8081"),
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT must be set")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
