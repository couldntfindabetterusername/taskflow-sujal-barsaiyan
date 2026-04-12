package config

import (
 "fmt"
 "os"
)

type Config struct {
 DatabaseURL string
 JWTSecret   string
 ServerPort  string
}

func Load() (*Config, error) {
 cfg := &Config{
 DatabaseURL: getEnv("DATABASE_URL", ""),
 JWTSecret:   getEnv("JWT_SECRET", ""),
 ServerPort:  getEnv("SERVER_PORT", "8080"),
 }

 if cfg.DatabaseURL == "" {
 return nil, fmt.Errorf("DATABASE_URL environment variable is required")
 }

 if cfg.JWTSecret == "" {
 return nil, fmt.Errorf("JWT_SECRET environment variable is required")
 }

 return cfg, nil
}

func getEnv(key, defaultValue string) string {
 value := os.Getenv(key)
 if value == "" {
 return defaultValue
 }
 return value
}
