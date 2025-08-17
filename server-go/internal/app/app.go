package app

import (
	"server/internal/api/story"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// Repo
	storyRepo := repository.NewInMemoryStoryRepo()

	// Service
	storySvc := service.NewStoryService(storyRepo)

	// Handler
	storyH := story.NewHandler(storySvc)

	// Routes
	api := r.Group("/api")
	v1 := api.Group("/v1")
	story.RegisterRoutes(v1, storyH)

	return r
}
