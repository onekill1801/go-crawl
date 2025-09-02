package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
	"time"
)

type ChapterService struct {
	repo repository.ChapterRepository
}

func NewChapterService(r repository.ChapterRepository) *ChapterService {
	return &ChapterService{repo: r}
}

func (s *ChapterService) Create(ctx context.Context, title, author, cover string) (*model.Chapter, error) {
	st := &model.Chapter{
		// ID:        uuid.NewString(),
		// Title:     title,
		// Author:    author,
		// CoverURL:  cover,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *ChapterService) Get(ctx context.Context, id string) (*model.Chapter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ChapterService) GetListByID(ctx context.Context, chapter_id, story_id string) ([]model.Image, error) {
	return s.repo.GetListByID(ctx, story_id, chapter_id)
}

func (s *ChapterService) List(ctx context.Context, offset, limit int) ([]model.Chapter, error) {
	return s.repo.List(ctx, offset, limit)
}
