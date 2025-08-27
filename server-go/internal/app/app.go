package app

import (
	"log"
	"server/internal/api/chapter"
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

	// Kết nối database
	mysqlDB, err := db.NewMySQL()
	if err != nil {
		log.Fatal(err)
	}

	// Khởi tạo các repository
	storyRepo := repository.NewMySQLStoryRepo(mysqlDB)
	if err != nil {
		log.Fatal(err)
	}
	chapterRepo := repository.NewMySQLChapterRepo(mysqlDB)
	if err != nil {
		log.Fatal(err)
	}
	// imagesRepo := repository.NewMySQLImageRepo(mysqlDB)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Khởi tạo các service và handler
	storySvc := service.NewStoryService(storyRepo)
	storyH := story.NewHandler(storySvc)

	chapterSvc := service.NewChapterService(chapterRepo)
	chapterH := chapter.NewHandler(chapterSvc)
	// imagesSvc := service.NewImageService(imagesRepo)
	// imageH := chapter.NewHandler(imagesSvc)

	// Thiết lập các route
	api := r.Group("/api")
	v1 := api.Group("/v1")
	story.RegisterRoutes(v1, storyH)
	chapter.RegisterRoutes(v1, chapterH)
	// story.RegisterRoutes(v1, imageH)

	return r, storyRepo
}
