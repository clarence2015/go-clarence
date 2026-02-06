package domain

import "time"

// PostStatus represents the publication status of a blog post.
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Post represents a blog article.
type Post struct {
	ID          uint       `gorm:"primaryKey"`
	Title       string     `gorm:"size:255;not null"`
	Slug        string     `gorm:"size:255;uniqueIndex;not null"`
	ContentMD   string     `gorm:"type:text;not null"`
	ContentHTML string     `gorm:"type:text"`
	Summary     string     `gorm:"type:text"`
	Status      PostStatus `gorm:"type:varchar(32);not null;default:'draft'"`
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Tag represents a tag that can be attached to posts.
type Tag struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:64;uniqueIndex;not null"`
	Slug      string    `gorm:"size:64;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Category represents a category that a post can belong to.
type Category struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:64;uniqueIndex;not null"`
	Slug      string    `gorm:"size:64;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

