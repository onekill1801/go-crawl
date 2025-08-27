package story

import (
	"net/http"
	"server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *service.StoryService }

func NewHandler(s *service.StoryService) *Handler { return &Handler{svc: s} }

func (h *Handler) Create(c *gin.Context) {
	var req CreateStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": err.Error()})
		return
	}
	st, err := h.svc.Create(c.Request.Context(), req.Title, req.Author, req.CoverURL)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, StoryResponse{
		ID: st.ID, Title: st.Title, Author: st.Author, CoverURL: st.CoverURL,
	})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	st, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, StoryResponse{
		ID: st.ID, Title: st.Title, Author: st.Author, CoverURL: st.CoverURL,
	})
}

func (h *Handler) List(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	// Convert string -> int
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	list, err := h.svc.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.Error(err)
		return
	}

	out := make([]StoryResponse, 0, len(list))
	for _, s := range list {
		out = append(out, StoryResponse{ID: s.ID, Title: s.Title, Author: s.Author, CoverURL: s.CoverURL, CreatedAt: s.CreatedAt})
	}
	c.JSON(http.StatusOK, out)
}
