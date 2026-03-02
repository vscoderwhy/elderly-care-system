package database

import (
	"context"
	"elderly-care-system/internal/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	return rdb
}
