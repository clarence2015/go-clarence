package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/clarence/GoBlog/internal/infra/markdown"
	"github.com/clarence/GoBlog/internal/infra/repository"
	"github.com/clarence/GoBlog/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter configures and returns a gin.Engine with basic routes.
func NewRouter(env string, logger *slog.Logger, db *gorm.DB) *gin.Engine {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(requestLogger(logger))

	// Wire up dependencies for use cases and handlers.
	postRepo := repository.NewPostRepository(db)
	mdRenderer := markdown.NewRenderer()
	postUC := usecase.NewPostUsecase(postRepo, mdRenderer)
	postHandlers := NewPostHandlers(postUC)

	// Health check endpoint.
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Public blog routes.
	engine.GET("/", postHandlers.ListRecent)
	engine.GET("/posts/:slug", postHandlers.GetBySlug)

	// Admin routes (JSON-based, authentication to be added later).
	admin := engine.Group("/admin")
	{
		admin.POST("/posts", postHandlers.Create)
	}

	return engine
}

// requestLogger logs basic HTTP request information using slog.
func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("http request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
		)
	}
}

