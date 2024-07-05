package core

import (
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func NewRedisLocker(redisClient redis.UniversalClient) *redislock.Client {
	locker := redislock.New(redisClient)
	return locker
}
