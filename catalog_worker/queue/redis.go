package queue

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(addr string) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisQueue{client: rdb}
}

func (q *RedisQueue) Push(queue string, data any) error {
	bytes, _ := json.Marshal(data)
	_, err := q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: queue,
		Values: map[string]interface{}{"data": string(bytes)},
	}).Result()
	return err
}
