package chapter

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/chapter")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetListImages)
	g.GET("/:id/next", h.GetListImagesNext)
	g.GET("/:id/previous", h.GetListImagesPrevious)
}
