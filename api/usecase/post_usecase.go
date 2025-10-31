package usecase

import (
	"context"
	"fmt"

	"github.com/rssh-jp/test-api/api/domain"
)

// PostUsecase defines business logic for posts
type PostUsecase interface {
	GetPosts(ctx context.Context, page, pageSize int) ([]domain.PostWithDetails, int64, error)
	GetPostByID(ctx context.Context, id int64) (*domain.PostWithDetails, error)
	GetPostBySlug(ctx context.Context, slug string) (*domain.PostWithDetails, error)
	GetPostsByCategory(ctx context.Context, categorySlug string, page, pageSize int) ([]domain.PostWithDetails, error)
	GetPostsByTag(ctx context.Context, tagSlug string, page, pageSize int) ([]domain.PostWithDetails, error)
	GetFeaturedPosts(ctx context.Context, limit int) ([]domain.PostWithDetails, error)
}

type postUsecase struct {
	postRepo domain.PostRepository
}

// NewPostUsecase creates a new post usecase
func NewPostUsecase(postRepo domain.PostRepository) PostUsecase {
	return &postUsecase{
		postRepo: postRepo,
	}
}

// GetPosts retrieves paginated posts
func (u *postUsecase) GetPosts(ctx context.Context, page, pageSize int) ([]domain.PostWithDetails, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, err := u.postRepo.FindAllWithDetails(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	total, err := u.postRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return posts, total, nil
}

// GetPostByID retrieves a post by ID and increments view count
func (u *postUsecase) GetPostByID(ctx context.Context, id int64) (*domain.PostWithDetails, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid post ID: %d", id)
	}

	post, err := u.postRepo.FindByIDWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Increment view count asynchronously (fire and forget)
	go func() {
		_ = u.postRepo.IncrementViewCount(context.Background(), id)
	}()

	return post, nil
}

// GetPostBySlug retrieves a post by slug and increments view count
func (u *postUsecase) GetPostBySlug(ctx context.Context, slug string) (*domain.PostWithDetails, error) {
	if slug == "" {
		return nil, fmt.Errorf("slug cannot be empty")
	}

	post, err := u.postRepo.FindBySlugWithDetails(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Increment view count asynchronously (fire and forget)
	go func() {
		_ = u.postRepo.IncrementViewCount(context.Background(), post.ID)
	}()

	return post, nil
}

// GetPostsByCategory retrieves posts by category
func (u *postUsecase) GetPostsByCategory(ctx context.Context, categorySlug string, page, pageSize int) ([]domain.PostWithDetails, error) {
	if categorySlug == "" {
		return nil, fmt.Errorf("category slug cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, err := u.postRepo.FindByCategoryWithDetails(ctx, categorySlug, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by category: %w", err)
	}

	return posts, nil
}

// GetPostsByTag retrieves posts by tag
func (u *postUsecase) GetPostsByTag(ctx context.Context, tagSlug string, page, pageSize int) ([]domain.PostWithDetails, error) {
	if tagSlug == "" {
		return nil, fmt.Errorf("tag slug cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, err := u.postRepo.FindByTagWithDetails(ctx, tagSlug, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by tag: %w", err)
	}

	return posts, nil
}

// GetFeaturedPosts retrieves featured posts
func (u *postUsecase) GetFeaturedPosts(ctx context.Context, limit int) ([]domain.PostWithDetails, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	posts, err := u.postRepo.FindFeaturedWithDetails(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured posts: %w", err)
	}

	return posts, nil
}
