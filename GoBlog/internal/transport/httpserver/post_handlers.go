package httpserver

import (
	"net/http"
	"strconv"

	"github.com/clarence/GoBlog/internal/usecase"
	"github.com/gin-gonic/gin"
)

// PostHandlers contains HTTP handlers for blog posts.
type PostHandlers struct {
	usecase *usecase.PostUsecase
}

// NewPostHandlers creates a new PostHandlers instance.
func NewPostHandlers(uc *usecase.PostUsecase) *PostHandlers {
	return &PostHandlers{usecase: uc}
}

// ListRecent handles GET / and returns a list of recent posts as JSON for now.
// Later this can be replaced with HTML templates.
func (h *PostHandlers) ListRecent(c *gin.Context) {
	ctx := c.Request.Context()

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	posts, err := h.usecase.ListRecent(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetBySlug handles GET /posts/:slug and returns a single post.
func (h *PostHandlers) GetBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	post, err := h.usecase.GetBySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Create handles POST /admin/posts for creating a new post.
// Expected JSON body:
// {
//   "title": "Post title",
//   "slug": "post-slug",
//   "content_md": "markdown content",
//   "summary": "short summary",
//   "publish": true
// }
func (h *PostHandlers) Create(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Slug      string `json:"slug" binding:"required"`
		ContentMD string `json:"content_md" binding:"required"`
		Summary   string `json:"summary"`
		Publish   bool   `json:"publish"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	post, err := h.usecase.CreatePost(ctx, usecase.CreatePostInput{
		Title:     req.Title,
		Slug:      req.Slug,
		ContentMD: req.ContentMD,
		Summary:   req.Summary,
		Publish:   req.Publish,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

