package ciudx

import (
	redis "github.com/dataspace-mobility/rs-iudx/ciudx/redis"
)

type App struct {
	// Redis connection
	RedisConnection *redis.RedisConnection
}

func NewApp() *App {
	return &App{
		RedisConnection: redis.NewRedisConnection(),
	}
}
