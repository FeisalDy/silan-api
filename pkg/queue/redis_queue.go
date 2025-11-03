package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisQueue handles publishing jobs to Redis queue
type RedisQueue struct {
	client    *redis.Client
	queueName string
}

// TranslationJobMessage represents the message structure sent to the queue
type TranslationJobMessage struct {
	JobID            string   `json:"job_id"`
	TargetLang       string   `json:"target_lang,omitempty"`
	SourceLang       string   `json:"source_lang,omitempty"`
	TargetFields     []string `json:"target_fields,omitempty"`
	EnableCodeFilter bool     `json:"enable_code_filter,omitempty"`
}

// NewRedisQueue creates a new Redis queue publisher
func NewRedisQueue(redisURL, queueName string) (*RedisQueue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisQueue{
		client:    client,
		queueName: queueName,
	}, nil
}

// Publish pushes a translation job to the Redis queue
func (q *RedisQueue) Publish(ctx context.Context, msg TranslationJobMessage) error {
	// Marshal message to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal job message: %w", err)
	}

	// Push to Redis list (queue)
	if err := q.client.RPush(ctx, q.queueName, data).Err(); err != nil {
		return fmt.Errorf("failed to push to Redis queue: %w", err)
	}

	return nil
}

// PublishBatch pushes multiple translation jobs to the Redis queue
func (q *RedisQueue) PublishBatch(ctx context.Context, messages []TranslationJobMessage) error {
	if len(messages) == 0 {
		return nil
	}

	// Prepare all messages
	values := make([]interface{}, len(messages))
	for i, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal job message at index %d: %w", i, err)
		}
		values[i] = data
	}

	// Push all messages in one operation
	if err := q.client.RPush(ctx, q.queueName, values...).Err(); err != nil {
		return fmt.Errorf("failed to push batch to Redis queue: %w", err)
	}

	return nil
}

// GetQueueLength returns the current number of jobs in the queue
func (q *RedisQueue) GetQueueLength(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, q.queueName).Result()
}

// Close closes the Redis connection
func (q *RedisQueue) Close() error {
	return q.client.Close()
}

// HealthCheck verifies the Redis connection is healthy
func (q *RedisQueue) HealthCheck(ctx context.Context) error {
	return q.client.Ping(ctx).Err()
}
