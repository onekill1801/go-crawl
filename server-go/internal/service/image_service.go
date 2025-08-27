package service

import (
	"context"
	"server/internal/model"
	"server/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ImageService struct {
	repo repository.ImageRepository
}

func NewImageService(r repository.ImageRepository) *ImageService {
	return &ImageService{repo: r}
}

func (s *ImageService) Create(ctx context.Context, title, author, cover string) (*model.Image, error) {
	st := &model.Image{
		ID:    uuid.NewString(),
		Title: title,
		// Author:    author,
		// CoverURL:  cover,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *ImageService) Get(ctx context.Context, id string) (*model.Image, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ImageService) List(ctx context.Context, offset, limit int) ([]model.Image, error) {
	return s.repo.List(ctx, offset, limit)
}
