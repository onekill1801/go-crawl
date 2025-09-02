package chapter

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/chapter")
	g.POST("", h.Create)
	g.GET("", h.List)                                               // Lấy danh sách chương
	g.GET("/:storyId/:chapterId", h.GetListImages)                  // Lấy danh danh image một chương theo ID
	g.GET("/:storyId/:chapterId/next", h.GetListImagesNext)         // Lấy danh sách image chương tiếp theo theo ID
	g.GET("/:storyId/:chapterId/previous", h.GetListImagesPrevious) // Lấy danh sách image chương trước theo ID
}
