package core

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURI string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr: opts.Addr,
	})
	return client, nil
}
