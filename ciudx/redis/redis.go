package redis

import (
	redis "github.com/go-redis/redis/v8"
)

// RedisConnection holds the redis connection.
type RedisConnection struct {
	client *redis.Client
}

// NewRedisConnection creates a new redis connection.
func NewRedisConnection() *RedisConnection {
	return &RedisConnection{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}
