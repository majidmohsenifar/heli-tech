package core

import (
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func NewRedisLocker(redisClient *redis.Client) *redislock.Client {
	locker := redislock.New(redisClient)
	return locker
}
