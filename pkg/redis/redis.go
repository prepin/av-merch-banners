package redis

import (
	"av-merch-shop/config"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedis(cfg config.RedisConfig, logger *slog.Logger) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis", "error", err)
	}

	return &Redis{client: client, logger: logger}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, json, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *Redis) Delete(ctx context.Context, key string) {
	err := r.client.Del(ctx, key)
	if err != nil {
		r.logger.Error("Failed to invalidate cache", "error", err)
	}
}
