package repository

import (
	"context"
	"encoding/json"

	"server/internal/model"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) GetItems(ctx context.Context) ([]model.Item, error) {
	res, err := r.client.XRange(ctx, "images_queue", "0", "+").Result()
	if err != nil {
		return nil, err
	}

	var items []model.Item
	for _, msg := range res {
		dataStr, _ := msg.Values["data"].(string)
		var parsed model.Item
		if err := json.Unmarshal([]byte(dataStr), &parsed); err == nil {
			items = append(items, parsed)
		}
	}
	return items, nil
}
