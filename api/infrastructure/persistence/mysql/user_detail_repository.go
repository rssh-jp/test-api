package mysql

import (
	"context"
	"database/sql"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rssh-jp/test-api/api/domain"
)

type userDetailRepository struct {
	db *sql.DB
}

// NewUserDetailRepository creates a new user detail repository
func NewUserDetailRepository(db *sql.DB) domain.UserDetailRepository {
	return &userDetailRepository{db: db}
}

func (r *userDetailRepository) FindDetailByID(ctx context.Context, id int64) (*domain.UserDetail, error) {
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := &newrelic.DatastoreSegment{
			StartTime:  txn.StartSegmentNow(),
			Product:    newrelic.DatastoreMySQL,
			Collection: "users",
			Operation:  "SELECT_WITH_JOINS",
		}
		defer segment.End()
	}

	detail := &domain.UserDetail{}

	// 1. ユーザー基本情報 + プロフィール情報を取得
	query := `
		SELECT 
			u.id, u.username, u.email, u.status, u.email_verified, u.last_login_at, u.created_at, u.updated_at,
			p.first_name, p.last_name, p.display_name, p.bio, p.avatar_url, p.birth_date, 
			p.gender, p.country_code, p.timezone, p.language_code, p.phone_number, p.website_url
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		WHERE u.id = ?
	`
	
	var profile domain.UserProfile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.ID, &detail.Username, &detail.Email, &detail.Status, &detail.EmailVerified, 
		&detail.LastLoginAt, &detail.CreatedAt, &detail.UpdatedAt,
		&profile.FirstName, &profile.LastName, &profile.DisplayName, &profile.Bio, &profile.AvatarURL,
		&profile.BirthDate, &profile.Gender, &profile.CountryCode, &profile.Timezone, 
		&profile.Language, &profile.PhoneNumber, &profile.WebsiteURL,
	)
	if err != nil {
		return nil, err
	}
	detail.Profile = &profile

	// 2. フォロー統計を取得
	followerQuery := `SELECT COUNT(*) FROM user_follows WHERE following_id = ?`
	followingQuery := `SELECT COUNT(*) FROM user_follows WHERE follower_id = ?`
	
	r.db.QueryRowContext(ctx, followerQuery, id).Scan(&detail.FollowStats.FollowerCount)
	r.db.QueryRowContext(ctx, followingQuery, id).Scan(&detail.FollowStats.FollowingCount)

	// 3. 投稿統計を取得
	statsQuery := `
		SELECT 
			COUNT(*) as post_count,
			COALESCE(SUM(view_count), 0) as total_views,
			COALESCE(SUM(like_count), 0) as total_likes
		FROM posts 
		WHERE user_id = ? AND status = 'published'
	`
	r.db.QueryRowContext(ctx, statsQuery, id).Scan(
		&detail.Stats.PostCount, 
		&detail.Stats.TotalViews, 
		&detail.Stats.TotalLikes,
	)

	// 4. コメント数を取得
	commentCountQuery := `SELECT COUNT(*) FROM comments WHERE user_id = ? AND status = 'approved'`
	r.db.QueryRowContext(ctx, commentCountQuery, id).Scan(&detail.Stats.CommentCount)

	// 5. 最近の投稿を取得（最新5件）
	postsQuery := `
		SELECT 
			id, title, slug, excerpt, status, published_at, 
			view_count, like_count, comment_count, is_featured, created_at
		FROM posts 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT 5
	`
	
	rows, err := r.db.QueryContext(ctx, postsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var post domain.UserPost
			rows.Scan(
				&post.ID, &post.Title, &post.Slug, &post.Excerpt, &post.Status, &post.PublishedAt,
				&post.ViewCount, &post.LikeCount, &post.CommentCount, &post.IsFeatured, &post.CreatedAt,
			)
			detail.RecentPosts = append(detail.RecentPosts, post)
		}
	}

	// 6. 最近のコメントを取得（最新5件）
	commentsQuery := `
		SELECT 
			c.id, c.post_id, p.title as post_title, c.content, c.status, c.like_count, c.created_at
		FROM comments c
		JOIN posts p ON c.post_id = p.id
		WHERE c.user_id = ?
		ORDER BY c.created_at DESC
		LIMIT 5
	`
	
	rows, err = r.db.QueryContext(ctx, commentsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var comment domain.UserComment
			rows.Scan(
				&comment.ID, &comment.PostID, &comment.PostTitle, &comment.Content, 
				&comment.Status, &comment.LikeCount, &comment.CreatedAt,
			)
			detail.RecentComments = append(detail.RecentComments, comment)
		}
	}

	// 7. 未読通知を取得（最新10件）
	notificationsQuery := `
		SELECT id, type, title, message, link_url, is_read, created_at, read_at
		FROM notifications
		WHERE user_id = ? AND is_read = FALSE
		ORDER BY created_at DESC
		LIMIT 10
	`
	
	rows, err = r.db.QueryContext(ctx, notificationsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var notification domain.UserNotification
			rows.Scan(
				&notification.ID, &notification.Type, &notification.Title, &notification.Message,
				&notification.LinkURL, &notification.IsRead, &notification.CreatedAt, &notification.ReadAt,
			)
			detail.UnreadNotifications = append(detail.UnreadNotifications, notification)
		}
	}

	// 8. スライスをnilから空配列に初期化
	if detail.RecentPosts == nil {
		detail.RecentPosts = []domain.UserPost{}
	}
	if detail.RecentComments == nil {
		detail.RecentComments = []domain.UserComment{}
	}
	if detail.UnreadNotifications == nil {
		detail.UnreadNotifications = []domain.UserNotification{}
	}

	return detail, nil
}

func (r *userDetailRepository) FindDetailByUsername(ctx context.Context, username string) (*domain.UserDetail, error) {
	// まずユーザー名からIDを取得
	var id int64
	query := `SELECT id FROM users WHERE username = ?`
	err := r.db.QueryRowContext(ctx, query, username).Scan(&id)
	if err != nil {
		return nil, err
	}
	
	// IDで詳細情報を取得
	return r.FindDetailByID(ctx, id)
}
