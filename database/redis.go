package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis() *redis.Client {
	addr := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return client
}
