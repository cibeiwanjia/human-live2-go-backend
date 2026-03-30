package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewClient(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
