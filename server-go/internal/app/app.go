package app

import (
	"log"
	"server/internal/api/story"
	"server/internal/db"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

// app/Setup.go
func Setup() (*gin.Engine, *repository.MySQLStoryRepo) {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	mysqlDB, err := db.NewMySQL()
	if err != nil {
		log.Fatal(err)
	}

	storyRepo := repository.NewMySQLStoryRepo(mysqlDB)
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
