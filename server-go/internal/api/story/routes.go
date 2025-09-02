package story

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/stories")
	g.POST("", h.Create)
	g.GET("", h.List)    // Lấy danh sách truyện
	g.GET("/:id", h.Get) // Lấy chi tiết một truyện theo ID
}
