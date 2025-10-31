package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rssh-jp/test-api/api/domain"
)

type postRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) domain.PostRepository {
	return &postRepository{db: db}
}

// FindAllWithDetails retrieves all published posts with joined data
func (r *postRepository) FindAllWithDetails(ctx context.Context, limit, offset int) ([]domain.PostWithDetails, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.category_id, p.title, p.slug, p.content, p.excerpt,
			p.status, p.published_at, p.view_count, p.like_count, p.comment_count,
			p.is_featured, p.created_at, p.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url,
			c.name as category_name,
			c.slug as category_slug
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.status = 'published' AND p.published_at IS NOT NULL
		ORDER BY p.published_at DESC
		LIMIT ? OFFSET ?
	`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "SELECT_WITH_JOIN",
		}
		defer segment.End()
	}

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.PostWithDetails
	postIDs := []int64{}

	for rows.Next() {
		var post domain.PostWithDetails
		err := rows.Scan(
			&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Slug,
			&post.Content, &post.Excerpt, &post.Status, &post.PublishedAt,
			&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.IsFeatured,
			&post.CreatedAt, &post.UpdatedAt,
			&post.AuthorUsername, &post.AuthorDisplayName, &post.AuthorAvatarURL,
			&post.CategoryName, &post.CategorySlug,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Load tags for all posts
	if len(postIDs) > 0 {
		tagsMap, err := r.loadTagsForPosts(ctx, postIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags: %w", err)
		}
		for i := range posts {
			if tags, ok := tagsMap[posts[i].ID]; ok {
				posts[i].Tags = tags
			}
		}
	}

	return posts, nil
}

// FindByIDWithDetails retrieves a post by ID with all related data
func (r *postRepository) FindByIDWithDetails(ctx context.Context, id int64) (*domain.PostWithDetails, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.category_id, p.title, p.slug, p.content, p.excerpt,
			p.status, p.published_at, p.view_count, p.like_count, p.comment_count,
			p.is_featured, p.created_at, p.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url,
			c.name as category_name,
			c.slug as category_slug
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ? AND p.status = 'published'
	`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "SELECT",
		}
		defer segment.End()
	}

	var post domain.PostWithDetails
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Slug,
		&post.Content, &post.Excerpt, &post.Status, &post.PublishedAt,
		&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.IsFeatured,
		&post.CreatedAt, &post.UpdatedAt,
		&post.AuthorUsername, &post.AuthorDisplayName, &post.AuthorAvatarURL,
		&post.CategoryName, &post.CategorySlug,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query post: %w", err)
	}

	// Load tags
	tags, err := r.loadTagsForPost(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	post.Tags = tags

	// Load latest comments
	comments, err := r.loadLatestCommentsForPost(ctx, post.ID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to load comments: %w", err)
	}
	post.LatestComments = comments

	return &post, nil
}

// FindBySlugWithDetails retrieves a post by slug with all related data
func (r *postRepository) FindBySlugWithDetails(ctx context.Context, slug string) (*domain.PostWithDetails, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.category_id, p.title, p.slug, p.content, p.excerpt,
			p.status, p.published_at, p.view_count, p.like_count, p.comment_count,
			p.is_featured, p.created_at, p.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url,
			c.name as category_name,
			c.slug as category_slug
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.slug = ? AND p.status = 'published'
	`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "SELECT",
		}
		defer segment.End()
	}

	var post domain.PostWithDetails
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Slug,
		&post.Content, &post.Excerpt, &post.Status, &post.PublishedAt,
		&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.IsFeatured,
		&post.CreatedAt, &post.UpdatedAt,
		&post.AuthorUsername, &post.AuthorDisplayName, &post.AuthorAvatarURL,
		&post.CategoryName, &post.CategorySlug,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query post: %w", err)
	}

	// Load tags
	tags, err := r.loadTagsForPost(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	post.Tags = tags

	// Load latest comments
	comments, err := r.loadLatestCommentsForPost(ctx, post.ID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to load comments: %w", err)
	}
	post.LatestComments = comments

	return &post, nil
}

// FindByCategoryWithDetails retrieves posts by category with related data
func (r *postRepository) FindByCategoryWithDetails(ctx context.Context, categorySlug string, limit, offset int) ([]domain.PostWithDetails, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.category_id, p.title, p.slug, p.content, p.excerpt,
			p.status, p.published_at, p.view_count, p.like_count, p.comment_count,
			p.is_featured, p.created_at, p.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url,
			c.name as category_name,
			c.slug as category_slug
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		INNER JOIN categories c ON p.category_id = c.id
		WHERE c.slug = ? AND p.status = 'published' AND p.published_at IS NOT NULL
		ORDER BY p.published_at DESC
		LIMIT ? OFFSET ?
	`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "SELECT_WITH_JOIN",
		}
		defer segment.End()
	}

	rows, err := r.db.QueryContext(ctx, query, categorySlug, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts by category: %w", err)
	}
	defer rows.Close()

	return r.scanPostsWithTags(ctx, rows)
}

// FindByTagWithDetails retrieves posts by tag with related data
func (r *postRepository) FindByTagWithDetails(ctx context.Context, tagSlug string, limit, offset int) ([]domain.PostWithDetails, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, p.category_id, p.title, p.slug, p.content, p.excerpt,
			p.status, p.published_at, p.view_count, p.like_count, p.comment_count,
			p.is_featured, p.created_at, p.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url,
			c.name as category_name,
			c.slug as category_slug
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN categories c ON p.category_id = c.id
		INNER JOIN post_tags pt ON p.id = pt.post_id
		INNER JOIN tags t ON pt.tag_id = t.id
		WHERE t.slug = ? AND p.status = 'published' AND p.published_at IS NOT NULL
		ORDER BY p.published_at DESC
		LIMIT ? OFFSET ?
	`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "SELECT_WITH_JOIN",
		}
		defer segment.End()
	}

	rows, err := r.db.QueryContext(ctx, query, tagSlug, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts by tag: %w", err)
	}
	defer rows.Close()

	return r.scanPostsWithTags(ctx, rows)
}

// FindFeaturedWithDetails retrieves featured posts with related data
func (r *postRepository) FindFeaturedWithDetails(ctx context.Context, limit int) ([]domain.PostWithDetails, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.category_id, p.title, p.slug, p.content, p.excerpt,
			p.status, p.published_at, p.view_count, p.like_count, p.comment_count,
			p.is_featured, p.created_at, p.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url,
			c.name as category_name,
			c.slug as category_slug
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.is_featured = TRUE AND p.status = 'published' AND p.published_at IS NOT NULL
		ORDER BY p.published_at DESC
		LIMIT ?
	`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "SELECT_WITH_JOIN",
		}
		defer segment.End()
	}

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query featured posts: %w", err)
	}
	defer rows.Close()

	return r.scanPostsWithTags(ctx, rows)
}

// GetTotalCount returns total count of published posts
func (r *postRepository) GetTotalCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM posts WHERE status = 'published' AND published_at IS NOT NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return count, nil
}

// IncrementViewCount increments the view count for a post
func (r *postRepository) IncrementViewCount(ctx context.Context, postID int64) error {
	query := `UPDATE posts SET view_count = view_count + 1 WHERE id = ?`

	// NewRelic automatically traces this query via context from nrecho middleware
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "posts",
			Operation:  "UPDATE",
		}
		defer segment.End()
	}

	_, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	return nil
}

// Helper function to load tags for multiple posts efficiently
func (r *postRepository) loadTagsForPosts(ctx context.Context, postIDs []int64) (map[int64][]domain.Tag, error) {
	if len(postIDs) == 0 {
		return make(map[int64][]domain.Tag), nil
	}

	// Build placeholders for IN clause
	placeholders := strings.Repeat("?,", len(postIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	query := fmt.Sprintf(`
		SELECT pt.post_id, t.id, t.name, t.slug, t.description, t.usage_count, t.created_at, t.updated_at
		FROM post_tags pt
		INNER JOIN tags t ON pt.tag_id = t.id
		WHERE pt.post_id IN (%s)
		ORDER BY t.name
	`, placeholders)

	// Convert postIDs to []interface{} for query
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	tagsMap := make(map[int64][]domain.Tag)
	for rows.Next() {
		var postID int64
		var tag domain.Tag
		err := rows.Scan(&postID, &tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.UsageCount, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tagsMap[postID] = append(tagsMap[postID], tag)
	}

	return tagsMap, rows.Err()
}

// Helper function to load tags for a single post
func (r *postRepository) loadTagsForPost(ctx context.Context, postID int64) ([]domain.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.usage_count, t.created_at, t.updated_at
		FROM post_tags pt
		INNER JOIN tags t ON pt.tag_id = t.id
		WHERE pt.post_id = ?
		ORDER BY t.name
	`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.UsageCount, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// Helper function to load latest comments for a post
func (r *postRepository) loadLatestCommentsForPost(ctx context.Context, postID int64, limit int) ([]domain.CommentWithAuthor, error) {
	query := `
		SELECT 
			c.id, c.post_id, c.user_id, c.parent_id, c.content, c.status,
			c.like_count, c.is_edited, c.created_at, c.updated_at,
			u.username as author_username,
			up.display_name as author_display_name,
			up.avatar_url as author_avatar_url
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		WHERE c.post_id = ? AND c.status = 'approved'
		ORDER BY c.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, postID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.CommentWithAuthor
	for rows.Next() {
		var comment domain.CommentWithAuthor
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID,
			&comment.Content, &comment.Status, &comment.LikeCount, &comment.IsEdited,
			&comment.CreatedAt, &comment.UpdatedAt,
			&comment.AuthorUsername, &comment.AuthorDisplayName, &comment.AuthorAvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// Helper function to scan posts and load their tags
func (r *postRepository) scanPostsWithTags(ctx context.Context, rows *sql.Rows) ([]domain.PostWithDetails, error) {
	var posts []domain.PostWithDetails
	postIDs := []int64{}

	for rows.Next() {
		var post domain.PostWithDetails
		err := rows.Scan(
			&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Slug,
			&post.Content, &post.Excerpt, &post.Status, &post.PublishedAt,
			&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.IsFeatured,
			&post.CreatedAt, &post.UpdatedAt,
			&post.AuthorUsername, &post.AuthorDisplayName, &post.AuthorAvatarURL,
			&post.CategoryName, &post.CategorySlug,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Load tags for all posts
	if len(postIDs) > 0 {
		tagsMap, err := r.loadTagsForPosts(ctx, postIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to load tags: %w", err)
		}
		for i := range posts {
			if tags, ok := tagsMap[posts[i].ID]; ok {
				posts[i].Tags = tags
			}
		}
	}

	return posts, nil
}
