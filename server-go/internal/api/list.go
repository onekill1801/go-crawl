package api

import (
	"net/http"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	service *service.ListService
}

func NewListHandler(s *service.ListService) *ListHandler {
	return &ListHandler{service: s}
}

func (h *ListHandler) GetList(c *gin.Context) {
	items, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
