package db

import (
	"crud/internal/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) *redis.Client {

	options, err := redis.ParseURL(cfg.RedisDSN)
	if err != nil {
		log.Fatalf("Failed to parse Redis DSN: %v", err)
	}

	rdb := redis.NewClient(options)

	return rdb

}
