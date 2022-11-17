package main

import (
	"os"

	"github.com/go-redis/redis/v9"
)

// Redis - Get the redis client
func Redis() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	client := redis.NewClient(opt)

	return client, err
}
