package chapter

import (
	"net/http"
	"server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *service.ChapterService }

func NewHandler(s *service.ChapterService) *Handler { return &Handler{svc: s} }

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
		ID: "st.ID", Title: st.Title, Author: "st.Author", CoverURL: "st.CoverURL",
	})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	st, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ChapterResponse{
		ID:        strconv.FormatInt(st.ID, 10),
		Title:     st.Title,
		StoryID:   st.StoryID,
		Content:   st.Content,
		ImageURL:  st.ImageURL,
		OrderStt:  st.OrderStt,
		CreatedAt: st.CreatedAt.String(),
	})
}

func (h *Handler) GetListImages(c *gin.Context) {
	id := c.Param("id")
	list, err := h.svc.GetListByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("X-Total-Count", strconv.Itoa(len(list)))

	out := make([]ChapterResponse, 0, len(list))
	for _, st := range list {
		out = append(out, ChapterResponse{
			ID:        strconv.FormatInt(st.ID, 10),
			Title:     st.Title,
			StoryID:   st.StoryID,
			Content:   st.Content,
			ImageURL:  st.ImageURL,
			OrderStt:  st.OrderStt,
			CreatedAt: st.CreatedAt.String(),
		})
	}
	c.JSON(http.StatusOK, out)

}

func (h *Handler) GetListImagesNext(c *gin.Context) {
	id := c.Param("id")
	st, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, StoryResponse{
		ID: "st.ID", Title: st.Title, Author: "st.Author", CoverURL: "st.CoverURL",
	})
}

func (h *Handler) GetListImagesPrevious(c *gin.Context) {
	id := c.Param("id")
	st, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, StoryResponse{
		ID: "st.ID", Title: st.Title, Author: "st.Author", CoverURL: "st.CoverURL",
	})
}

func (h *Handler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context(), 0, 50)
	if err != nil {
		c.Error(err)
		return
	}

	out := make([]StoryResponse, 0, len(list))
	for _, s := range list {
		out = append(out, StoryResponse{ID: "s.ID", Title: s.Title, Author: "s.Author", CoverURL: "s.CoverURL"})
	}
	c.JSON(http.StatusOK, out)
}
