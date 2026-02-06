package repository

import (
	"context"

	"github.com/clarence/GoBlog/internal/domain"
	"gorm.io/gorm"
)

// PostRepository defines persistence operations for posts.
type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id uint) (*domain.Post, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Post, error)
	ListRecent(ctx context.Context, limit int) ([]domain.Post, error)
}

type postGormRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new GORM-based PostRepository implementation.
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postGormRepository{db: db}
}

func (r *postGormRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postGormRepository) FindByID(ctx context.Context, id uint) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postGormRepository) FindBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postGormRepository) ListRecent(ctx context.Context, limit int) ([]domain.Post, error) {
	if limit <= 0 {
		limit = 10
	}

	var posts []domain.Post
	if err := r.db.WithContext(ctx).
		Order("published_at DESC NULLS LAST, created_at DESC").
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}


