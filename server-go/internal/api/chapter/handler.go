package chapter

import (
	"fmt"
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

func (h *Handler) getListImages(c *gin.Context, delta int) {
	chapterIDStr := c.Param("chapterId")
	chapterID, err := strconv.Atoi(chapterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapterId"})
		return
	}

	storyID := c.Param("storyId")
	targetChapterID := strconv.Itoa(chapterID + delta)

	// debug query param (offset, limit)
	offset := c.Query("offset")
	limit := c.Query("limit")
	fmt.Println(">>> offset:", offset, " - limit:", limit)
	fmt.Println(">>> chapterId:", targetChapterID, " - storyId:", storyID)

	// query DB
	list, err := h.svc.GetListByID(c.Request.Context(), storyID, targetChapterID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(len(list)))

	out := make([]ImageResponse, 0, len(list))
	for _, st := range list {
		out = append(out, ImageResponse{
			ID:        st.ID,
			ChapterId: st.ChapterId,
			StoryID:   st.StoryID,
			ImageURL:  st.ImageURL,
			OrderStt:  st.OrderStt,
			Referer:   st.Referer,
			CreatedAt: st.CreatedAt.String(),
		})
	}
	c.JSON(http.StatusOK, out)
}

// === Handler public ===

func (h *Handler) GetListImages(c *gin.Context) {
	h.getListImages(c, 0) // giữ nguyên chapterId
}

func (h *Handler) GetListImagesNext(c *gin.Context) {
	h.getListImages(c, 1) // chapterId +1
}

func (h *Handler) GetListImagesPrevious(c *gin.Context) {
	h.getListImages(c, -1) // chapterId -1
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
