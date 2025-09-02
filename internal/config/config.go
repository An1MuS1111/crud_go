package config

import "os"

type Config struct {
	PostgresDSN string
	RedisDSN    string
}

func Load() *Config {
	return &Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@db/postgres?sslmode=disable"),
		RedisDSN:    getEnv("REDIS_DSN", "redis://:password@redis:6379/0"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
