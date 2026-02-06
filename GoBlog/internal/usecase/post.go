package usecase

import (
	"context"
	"time"

	"github.com/clarence/GoBlog/internal/domain"
)

// PostRepository defines the persistence operations required by the post use case.
// It is satisfied by the infra/repository implementation.
type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id uint) (*domain.Post, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Post, error)
	ListRecent(ctx context.Context, limit int) ([]domain.Post, error)
}

// MarkdownRenderer converts markdown content into HTML.
type MarkdownRenderer interface {
	RenderToHTML(markdown string) (string, error)
}

// PostUsecase contains business logic for blog posts.
type PostUsecase struct {
	repo     PostRepository
	renderer MarkdownRenderer
}

// NewPostUsecase creates a new PostUsecase.
func NewPostUsecase(repo PostRepository, renderer MarkdownRenderer) *PostUsecase {
	return &PostUsecase{
		repo:     repo,
		renderer: renderer,
	}
}

// CreatePostInput represents data required to create a post.
type CreatePostInput struct {
	Title     string
	Slug      string
	ContentMD string
	Summary   string
	Publish   bool
}

// CreatePost creates a new post, rendering markdown to HTML and setting status.
func (u *PostUsecase) CreatePost(ctx context.Context, in CreatePostInput) (*domain.Post, error) {
	html, err := u.renderer.RenderToHTML(in.ContentMD)
	if err != nil {
		return nil, err
	}

	post := &domain.Post{
		Title:       in.Title,
		Slug:        in.Slug,
		ContentMD:   in.ContentMD,
		ContentHTML: html,
		Summary:     in.Summary,
		Status:      domain.PostStatusDraft,
	}

	if in.Publish {
		now := time.Now()
		post.Status = domain.PostStatusPublished
		post.PublishedAt = &now
	}

	if err := u.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetBySlug retrieves a single post by its slug.
func (u *PostUsecase) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	return u.repo.FindBySlug(ctx, slug)
}

// ListRecent returns recent posts ordered by published/created time.
func (u *PostUsecase) ListRecent(ctx context.Context, limit int) ([]domain.Post, error) {
	return u.repo.ListRecent(ctx, limit)
}

