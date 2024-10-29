package redis

import (
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/redis/go-redis/v9"
	"sync"
)

var once sync.Once
var redisClient *redis.Client

func Connect(config *app.Configuration) *redis.Client {
	//var redisClient *redis.Client
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})
		defer func(redisClient *redis.Client) {
			err := redisClient.Close()
			if err != nil {
				wrapper.LogError(fmt.Sprintf("Error closing redis connection: %v", err),
					"main")
			}
		}(redisClient)
	})

	return redisClient
}
