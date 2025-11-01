package domain

import (
	"context"
	"time"
)

// UserProfile はユーザープロフィール詳細情報
type UserProfile struct {
	FirstName   *string `json:"firstName,omitempty"`
	LastName    *string `json:"lastName,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	BirthDate   *string `json:"birthDate,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	CountryCode *string `json:"countryCode,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Language    *string `json:"language,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	WebsiteURL  *string `json:"websiteUrl,omitempty"`
}

// UserPost はユーザーの投稿情報
type UserPost struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Excerpt     *string   `json:"excerpt,omitempty"`
	Status      string    `json:"status"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	ViewCount   int       `json:"viewCount"`
	LikeCount   int       `json:"likeCount"`
	CommentCount int      `json:"commentCount"`
	IsFeatured  bool      `json:"isFeatured"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserComment はユーザーのコメント情報
type UserComment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	PostTitle string    `json:"postTitle"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	LikeCount int       `json:"likeCount"`
	CreatedAt time.Time `json:"createdAt"`
}

// FollowStats はフォロー統計情報
type FollowStats struct {
	FollowerCount  int `json:"followerCount"`
	FollowingCount int `json:"followingCount"`
}

// UserStats はユーザー統計情報
type UserStats struct {
	PostCount    int `json:"postCount"`
	CommentCount int `json:"commentCount"`
	TotalLikes   int `json:"totalLikes"`
	TotalViews   int `json:"totalViews"`
}

// UserNotification はユーザーの通知情報
type UserNotification struct {
	ID        int64      `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	LinkURL   *string    `json:"linkUrl,omitempty"`
	IsRead    bool       `json:"isRead"`
	CreatedAt time.Time  `json:"createdAt"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
}

// UserDetail は全ての関連情報を含むユーザー詳細
type UserDetail struct {
	// 基本情報
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	EmailVerified  bool      `json:"emailVerified"`
	LastLoginAt    *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	
	// プロフィール情報
	Profile *UserProfile `json:"profile,omitempty"`
	
	// フォロー情報
	FollowStats FollowStats `json:"followStats"`
	
	// 統計情報
	Stats UserStats `json:"stats"`
	
	// 最近の投稿（最新5件）
	RecentPosts []UserPost `json:"recentPosts"`
	
	// 最近のコメント（最新5件）
	RecentComments []UserComment `json:"recentComments"`
	
	// 未読通知（最新10件）
	UnreadNotifications []UserNotification `json:"unreadNotifications"`
}

// UserDetailRepository はユーザー詳細情報のリポジトリインターフェース
type UserDetailRepository interface {
	FindDetailByID(ctx context.Context, id int64) (*UserDetail, error)
	FindDetailByUsername(ctx context.Context, username string) (*UserDetail, error)
}
