package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tharun-rs/rprox/config"
)

type RedisClient struct {
	rdb *redis.Client
	ctx context.Context
}

func (c *RedisClient) Init(conf config.Config) error {
	c.ctx = context.Background()
	c.rdb = redis.NewClient(&redis.Options{
		Addr:     conf.RedisURL,
		Password: conf.RedisPass,
		DB:       conf.RedisDB,
	})

	_, err := c.rdb.Ping(c.ctx).Result()
	return err
}

func (c *RedisClient) Put(key string, value string, expiration time.Duration) error {
	err := c.rdb.Set(c.ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}
	return nil
}

func (c *RedisClient) Get(key string) (string, error) {
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}
	return val, nil
}

func (c *RedisClient) Delete(key string) error {
	err := c.rdb.Del(c.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

func (c *RedisClient) Close() error {
	err := c.rdb.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	return nil
}

func (c *RedisClient) Extend(key string, expiration time.Duration) error {
	err := c.rdb.Expire(c.ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to extend expiration: %w", err)
	}
	return nil
}
