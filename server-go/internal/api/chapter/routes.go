package story

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/chapter")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.GET("/:id/next", h.Get)
	g.GET("/:id/previous", h.Get)
}
