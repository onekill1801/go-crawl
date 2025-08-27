package redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

// NewClient khởi tạo Redis client
func NewClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

// ReadStream đọc message từ Redis Stream
func ReadStream(ctx context.Context, rdb *redis.Client, stream, group, consumer string) ([]redis.XMessage, error) {
	res, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
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

func EnsureConsumerGroup(rdb *redis.Client, stream, group string, reset bool) error {
	ctx := context.Background()

	if reset {
		// Xóa group nếu có
		_ = rdb.XGroupDestroy(ctx, stream, group).Err()
		// Không cần check error ở đây, nếu group chưa có sẽ trả lỗi nhưng ta bỏ qua
	}

	// Tạo lại group
	startID := "$" // chỉ đọc message mới
	if reset {
		startID = "0" // đọc lại toàn bộ từ đầu
	}

	err := rdb.XGroupCreateMkStream(ctx, stream, group, startID).Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}
