package cache

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/pkg/redis"
	"context"
	"fmt"
	"time"
)

type RedisUserInfoCache struct {
	redis *redis.Redis
}

func NewUserInfoCache(redis *redis.Redis) *RedisUserInfoCache {
	return &RedisUserInfoCache{
		redis: redis,
	}
}

func (c *RedisUserInfoCache) getUserInfoKey(userID int) string {
	return fmt.Sprintf("user_info:%d", userID)
}

func (c *RedisUserInfoCache) GetUserInfo(ctx context.Context, userID int) (*entities.UserInfo, error) {
	var info entities.UserInfo
	err := c.redis.Get(ctx, c.getUserInfoKey(userID), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *RedisUserInfoCache) SetUserInfo(ctx context.Context, userID int, info entities.UserInfo) error {
	return c.redis.Set(ctx, c.getUserInfoKey(userID), &info, 15*time.Second)
}

func (c *RedisUserInfoCache) ExpireUserInfo(ctx context.Context, userID int) {
	c.redis.Delete(ctx, c.getUserInfoKey(userID))
}
