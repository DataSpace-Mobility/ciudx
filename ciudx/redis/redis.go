package redis

import (
	utils "github.com/dataspace-mobility/rs-iudx/ciudx/utils"
	redis "github.com/go-redis/redis/v8"
)

// Connection holds the redis connection.
type Connection struct {
	Client *redis.Client
}

// NewRedisConnection creates a new redis connection.
func NewRedisConnection() *Connection {
	host := utils.Getenv("REDIS_HOST", "ds_redis")
	port := utils.Getenv("REDIS_PORT", "6379")
	username := utils.Getenv("REDIS_USERNAME", "")
	password := utils.Getenv("REDIS_PASSWORD", "")

	return &Connection{
		Client: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Username: username,
			Password: password,
			DB:       0,
		}),
	}
}
