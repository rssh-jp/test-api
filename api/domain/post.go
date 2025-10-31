package domain

import (
	"context"
	"time"
)

// Post represents a blog post entity
type Post struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"userId"`
	CategoryID    *int64    `json:"categoryId,omitempty"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Content       string    `json:"content"`
	Excerpt       *string   `json:"excerpt,omitempty"`
	Status        string    `json:"status"`
	PublishedAt   *time.Time `json:"publishedAt,omitempty"`
	ViewCount     int32     `json:"viewCount"`
	LikeCount     int32     `json:"likeCount"`
	CommentCount  int32     `json:"commentCount"`
	IsFeatured    bool      `json:"isFeatured"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PostWithDetails represents a post with joined related data
type PostWithDetails struct {
	Post
	// Author information
	AuthorUsername    string  `json:"authorUsername"`
	AuthorDisplayName *string `json:"authorDisplayName,omitempty"`
	AuthorAvatarURL   *string `json:"authorAvatarUrl,omitempty"`
	
	// Category information
	CategoryName *string `json:"categoryName,omitempty"`
	CategorySlug *string `json:"categorySlug,omitempty"`
	
	// Tags
	Tags []Tag `json:"tags,omitempty"`
	
	// Comment preview (latest comments)
	LatestComments []CommentWithAuthor `json:"latestComments,omitempty"`
}

// Category represents a post category
type Category struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	ParentID     *int64    `json:"parentId,omitempty"`
	DisplayOrder int32     `json:"displayOrder"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Tag represents a post tag
type Tag struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	UsageCount  int32     `json:"usageCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	UserID    int64     `json:"userId"`
	ParentID  *int64    `json:"parentId,omitempty"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	LikeCount int32     `json:"likeCount"`
	IsEdited  bool      `json:"isEdited"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CommentWithAuthor represents a comment with author information
type CommentWithAuthor struct {
	Comment
	AuthorUsername    string  `json:"authorUsername"`
	AuthorDisplayName *string `json:"authorDisplayName,omitempty"`
	AuthorAvatarURL   *string `json:"authorAvatarUrl,omitempty"`
}

// PostRepository defines methods for post data access
type PostRepository interface {
	// FindAllWithDetails retrieves all published posts with joined data
	FindAllWithDetails(ctx context.Context, limit, offset int) ([]PostWithDetails, error)
	
	// FindByIDWithDetails retrieves a post by ID with all related data
	FindByIDWithDetails(ctx context.Context, id int64) (*PostWithDetails, error)
	
	// FindBySlugWithDetails retrieves a post by slug with all related data
	FindBySlugWithDetails(ctx context.Context, slug string) (*PostWithDetails, error)
	
	// FindByCategoryWithDetails retrieves posts by category with related data
	FindByCategoryWithDetails(ctx context.Context, categorySlug string, limit, offset int) ([]PostWithDetails, error)
	
	// FindByTagWithDetails retrieves posts by tag with related data
	FindByTagWithDetails(ctx context.Context, tagSlug string, limit, offset int) ([]PostWithDetails, error)
	
	// FindFeaturedWithDetails retrieves featured posts with related data
	FindFeaturedWithDetails(ctx context.Context, limit int) ([]PostWithDetails, error)
	
	// GetTotalCount returns total count of published posts
	GetTotalCount(ctx context.Context) (int64, error)
	
	// IncrementViewCount increments the view count for a post
	IncrementViewCount(ctx context.Context, postID int64) error
}
