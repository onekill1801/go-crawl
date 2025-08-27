package repository

import (
	"context"
	"server/internal/model"
)

type StoryRepository interface {
	Create(ctx context.Context, s *model.Story) error
	GetByID(ctx context.Context, id string) (*model.Story, error)
	List(ctx context.Context, offset, limit int) ([]model.Story, error)
}

type ChapterRepository interface {
	Create(ctx context.Context, s *model.Chapter) error
	GetByID(ctx context.Context, id string) (*model.Chapter, error)
	List(ctx context.Context, offset, limit int) ([]model.Chapter, error)
	ListImages(ctx context.Context, offset, limit int) ([]model.Chapter, error)
	ListImagesNext(ctx context.Context, offset, limit int) ([]model.Chapter, error)
	ListImagesPrevious(ctx context.Context, offset, limit int) ([]model.Chapter, error)
}

type ImageRepository interface {
	Create(ctx context.Context, s *model.Image) error
	GetByID(ctx context.Context, id string) (*model.Image, error)
	List(ctx context.Context, offset, limit int) ([]model.Image, error)
}

type InMemoryStoryRepo struct {
	data map[string]model.Story
}

func NewInMemoryStoryRepo() *InMemoryStoryRepo {
	return &InMemoryStoryRepo{data: map[string]model.Story{}}
}
func (r *InMemoryStoryRepo) Create(_ context.Context, s *model.Story) error {
	r.data[s.ID] = *s
	return nil
}
func (r *InMemoryStoryRepo) GetByID(_ context.Context, id string) (*model.Story, error) {
	s, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &s, nil
}
func (r *InMemoryStoryRepo) List(_ context.Context, offset, limit int) ([]model.Story, error) {
	out := make([]model.Story, 0, len(r.data))
	for _, v := range r.data {
		out = append(out, v)
	}
	return out, nil
}
