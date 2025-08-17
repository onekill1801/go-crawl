package main

import (
	"server/internal/api/story"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// Repos (ở đây dùng in-memory; sau đổi MySQL/Redis)
	storyRepo := repository.NewInMemoryStoryRepo()

	// Services
	storySvc := service.NewStoryService(storyRepo)

	// Handlers
	storyH := story.NewHandler(storySvc)

	// Versioning + route groups
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")
	story.RegisterRoutes(v1, storyH)
	story.RegisterRoutes(v2, storyH)

	r.Run(":8080")
}
