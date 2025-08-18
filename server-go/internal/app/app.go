package app

import (
	"log"
	"server/internal/api/story"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

// app/Setup.go
func Setup() (*gin.Engine, *repository.MySQLStoryRepo) {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	storyRepo, err := repository.NewMySQLStoryRepo()
	if err != nil {
		log.Fatal(err)
	}

	storySvc := service.NewStoryService(storyRepo)
	storyH := story.NewHandler(storySvc)

	api := r.Group("/api")
	v1 := api.Group("/v1")
	story.RegisterRoutes(v1, storyH)

	return r, storyRepo
}
