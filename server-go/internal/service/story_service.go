package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
	"time"

	"github.com/google/uuid"
)

type StoryService struct {
	repo repository.StoryRepository
}

func NewStoryService(r repository.StoryRepository) *StoryService {
	return &StoryService{repo: r}
}

func (s *StoryService) Create(ctx context.Context, title, author, cover string) (*model.Story, error) {
	st := &model.Story{
		ID:        uuid.NewString(),
		Title:     title,
		Author:    author,
		CoverURL:  cover,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *StoryService) Get(ctx context.Context, id string) (*model.Story, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *StoryService) List(ctx context.Context, offset, limit int) ([]model.Story, error) {
	return s.repo.List(ctx, offset, limit)
}
