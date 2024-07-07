package core

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURI string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})
	return client, nil
}
