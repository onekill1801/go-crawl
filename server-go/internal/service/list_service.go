package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
)

type ListService struct {
	repo *repository.RedisRepo
}

func NewListService(r *repository.RedisRepo) *ListService {
	return &ListService{repo: r}
}

func (s *ListService) GetAll(ctx context.Context) ([]model.Item, error) {
	items, err := s.repo.GetItems(ctx)
	if err != nil {
		return nil, err
	}
	// business logic thêm (vd: append 1 item cứng)
	items = append(items, model.Item{
		Referer:  "https://ac.qq.com/",
		ImageURL: "https://example.com/fallback.jpg",
		Title:    "Thêm tay",
		Order:    999,
	})
	return items, nil
}
