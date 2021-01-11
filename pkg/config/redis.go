package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// ConnectRedis creates a redis context
func ConnectRedis(ctx context.Context, env Env) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", env.RedisHost, env.RedisPort),
		Password: env.RedisPassword,
		DB:       0, // use default DB
	})

	_, err := client.Ping(ctx).Result()

	return client, err
}
