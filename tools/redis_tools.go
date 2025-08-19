package tools

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisClient *redis.Client

// InitRedisClient initializes the Redis client.
func InitRedisClient() {
	var config, _ = GetAppConfig("config.yaml")
	redisClient = redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        config.Redis.Addr + ":" + config.Redis.Port,
		Password:    config.Redis.Password, // Specify your Redis password if required
		DB:          0,                     // Specify your Redis database number
		MaxRetries:  3,
		DialTimeout: 30 * time.Second,
		PoolSize:    10,
	})
}

// GetRedisClient returns the Redis client.
func Redis_GetRedisClient() *redis.Client {
	if redisClient == nil {
		InitRedisClient()
	}
	return redisClient
}

// GetValue retrieves the value for the specified key from Redis.
func Redis_GetValue(key string) (string, error) {
	ctx := context.Background()
	result, err := Redis_GetRedisClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("Key '%s' does not exist", key)
	} else if err != nil {
		return "", err
	}
	return result, nil
}

// SetValue sets a key-value pair in Redis with an optional expiration time.
func Redis_SetValue(key, value string, expiration time.Duration) (error, error) {
	ctx := context.Background()
	err := Redis_GetRedisClient().Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err, nil
	}
	return nil, nil
}

// SetValue sets a key-value pair in Redis with an optional expiration time.
func Redis_SetchatMessagesValue(key string, value string, expiration time.Duration) (error, error) {
	ctx := context.Background()
	err := Redis_GetRedisClient().Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err, nil
	}
	return nil, nil
}

// SetValue sets a key-value pair in Redis with an optional expiration time.
func Redis_GetchatMessagesValue(key string) string {
	ctx := context.Background()
	result, _ := Redis_GetRedisClient().Get(ctx, key).Result()

	return result
}

// HashGetValue sets a key-value pair in Redis with an optional expiration time.
func Redis_HashGetValue(hashkey, key string) string {
	ctx := context.Background()
	hSet := Redis_GetRedisClient().HGet(ctx, hashkey, key)
	return hSet.Val()
}

// HashSetValue sets a key-value pair in Redis with an optional expiration time.
func Redis_HashSetValue(hashkey, key, value string) *redis.IntCmd {
	ctx := context.Background()
	hSet := Redis_GetRedisClient().HSet(ctx, hashkey, key, value)
	return hSet
}

// HashSetValue sets a key-value pair in Redis with an optional expiration time.
func Redis_ExistsValue(key string) int64 {
	ctx := context.Background()
	exit_status, _ := Redis_GetRedisClient().Exists(ctx, key).Result()
	return exit_status
}

// Other utility functions for setting, updating, or deleting values can be added here.
