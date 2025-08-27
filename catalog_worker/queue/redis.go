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

func (q *RedisQueue) ReadStream(ctx context.Context, stream, group, consumer string) ([]redis.XMessage, error) {
	res, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    10,
		Block:    0, // block cho đến khi có message
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return res[0].Messages, nil
	}
	return nil, nil
}

func (q *RedisQueue) Push(queue string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var values map[string]interface{}
	if err := json.Unmarshal(bytes, &values); err != nil {
		return err
	}

	_, err = q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: queue,
		Values: values,
	}).Result()

	return err
}
