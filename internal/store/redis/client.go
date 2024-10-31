package redis

import (
	"github.com/redis/go-redis/v9"
)

func Connect(addr, password string, db int) *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return redisClient
}
