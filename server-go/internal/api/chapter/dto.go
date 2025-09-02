package chapter

type CreateStoryRequest struct {
	Title    string `json:"title"  binding:"required,min=1"`
	Author   string `json:"author" binding:"required"`
	CoverURL string `json:"cover_url" binding:"omitempty,url"`
}

type StoryResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	CoverURL string `json:"cover_url"`
}

type CreateChapterRequest struct {
	Title    string `json:"title"  binding:"required,min=1"`
	Author   string `json:"author" binding:"required"`
	CoverURL string `json:"cover_url" binding:"omitempty,url"`
}

type ChapterResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StoryID   string `json:"story_id"`
	Content   string `json:"content"`
	ImageURL  string `json:"image_url"`
	OrderStt  int64  `json:"order_stt"`
	CreatedAt string `json:"created_at"`
}

type CreateImageRequest struct {
	Title    string `json:"title"  binding:"required,min=1"`
	Author   string `json:"author" binding:"required"`
	CoverURL string `json:"cover_url" binding:"omitempty,url"`
}

type ImageResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StoryID   string `json:"story_id"`
	ImageURL  string `json:"image_url"`
	OrderStt  int64  `json:"order_stt"`
	CreatedAt string `json:"created_at"`
	ChapterId int64  `json:"chapter_id"`
	Referer   string `json:"referer"`
}
