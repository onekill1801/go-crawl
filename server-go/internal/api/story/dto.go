package story

import "time"

type CreateStoryRequest struct {
	Title    string `json:"title"  binding:"required,min=1"`
	Author   string `json:"author" binding:"required"`
	CoverURL string `json:"cover_url" binding:"omitempty,url"`
}

type StoryResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CoverURL  string    `json:"cover_url"`
	CreatedAt time.Time `json:"created_at"`
}
