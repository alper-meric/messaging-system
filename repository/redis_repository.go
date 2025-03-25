package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRepository implements the CacheRepository interface using Redis
type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisRepository creates a new RedisRepository instance
func NewRedisRepository(addr, password string, db int) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Check connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis server: %v", err)
	}

	return &RedisRepository{
		client: client,
		ctx:    ctx,
	}, nil
}

// CacheMessageID saves a message ID and send time to Redis
func (r *RedisRepository) CacheMessageID(messageID string, sentAt time.Time) error {
	// Save message ID as key and send time as value
	key := fmt.Sprintf("message:%s", messageID)
	err := r.client.Set(r.ctx, key, sentAt.Format(time.RFC3339), 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("redis cache error: %v", err)
	}
	return nil
}

// GetCachedMessage retrieves message information from Redis
func (r *RedisRepository) GetCachedMessage(messageID string) (time.Time, error) {
	key := fmt.Sprintf("message:%s", messageID)
	result, err := r.client.Get(r.ctx, key).Result()

	if err == redis.Nil {
		return time.Time{}, fmt.Errorf("message ID not found: %s", messageID)
	} else if err != nil {
		return time.Time{}, fmt.Errorf("redis read error: %v", err)
	}

	// Parse time information
	sentAt, err := time.Parse(time.RFC3339, result)
	if err != nil {
		return time.Time{}, fmt.Errorf("time format error: %v", err)
	}

	return sentAt, nil
}
