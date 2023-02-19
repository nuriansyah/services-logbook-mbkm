package config

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(config Config) (*redis.Client, error) {
	host := config.Get("REDIS_HOST")
	port := config.Get("REDIS_PORT")
	password := config.Get("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
func NewInitializedRedis(config Config) (*redis.Client, error) {
	client, err := NewRedisClient(config)
	if err != nil {
		return nil, err
	}
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
