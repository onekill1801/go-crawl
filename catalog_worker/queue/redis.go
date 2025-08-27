package queue

import (
	"context"
	"encoding/json"
	"strings"
	"time"

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

func (q *RedisQueue) EnsureConsumerGroup(stream, group string, reset bool) error {
	ctx := context.Background()

	if reset {
		_ = q.client.XGroupDestroy(ctx, stream, group).Err()
	}

	// Tạo lại group
	startID := "$" // chỉ đọc message mới
	if reset {
		startID = "0" // đọc lại toàn bộ từ đầu
	}

	err := q.client.XGroupCreateMkStream(ctx, stream, group, startID).Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

func (q *RedisQueue) ReadStream(ctx context.Context, stream, group, consumer string) ([]redis.XMessage, error) {
	res, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    1,
		Block:    0, // block cho đến khi có message
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(res) > 0 && len(res[0].Messages) > 0 {
		time.Sleep(1 * time.Second)
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
