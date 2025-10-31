package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rssh-jp/test-api/api/domain"
)

// cachedUserRepository はキャッシュのためのDecorator/Proxyパターンを実装します。
// 任意のdomain.UserRepository実装（MySQL, PostgreSQLなど）をラップし、
// Redisキャッシュ機能を追加します。
//
// クリーンアーキテクチャ的に正しい理由:
// - domain.UserRepositoryインターフェースに依存（具体実装ではない）
// - このクラスとラップされるリポジトリは両方ともInfrastructure層
// - 依存関係の方向が内側を向いている（Infrastructure -> Domain）
type cachedUserRepository struct {
	baseRepo    domain.UserRepository // Domainインターフェース - 任意の実装が可能
	redisClient *redis.Client
	ctx         context.Context
	ttl         time.Duration
}

// NewCachedUserRepository は新しいキャッシュ付きユーザーリポジトリを作成します。
// baseRepoはdomain.UserRepositoryの任意の実装（MySQL, PostgreSQLなど）が使用できます。
// Decoratorパターンに従い、透過的にキャッシュ機能を追加します。
func NewCachedUserRepository(baseRepo domain.UserRepository, redisClient *redis.Client) domain.UserRepository {
	return &cachedUserRepository{
		baseRepo:    baseRepo,
		redisClient: redisClient,
		ctx:         context.Background(),
		ttl:         5 * time.Minute, // Cache TTL: 5 minutes
	}
}

func (r *cachedUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	cacheKey := "users:all"

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []domain.User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			log.Printf("✓ Redis Cache HIT: %s", cacheKey)
			return users, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL", cacheKey)
	users, err := r.baseRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(users)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return users, nil
}

func (r *cachedUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	cacheKey := getCacheKey(id)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			log.Printf("✓ Redis Cache HIT: %s", cacheKey)
			return &user, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL", cacheKey)
	user, err := r.baseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(user)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return user, nil
}

func (r *cachedUserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.baseRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	// Invalidate list cache
	r.redisClient.Del(ctx, "users:all")
	log.Printf("⚠ Redis Cache INVALIDATE: users:all (User created: ID=%d)", user.ID)

	return nil
}

func (r *cachedUserRepository) Update(ctx context.Context, user *domain.User) error {
	err := r.baseRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Invalidate caches
	r.redisClient.Del(ctx, getCacheKey(user.ID))
	r.redisClient.Del(ctx, "users:all")
	log.Printf("⚠ Redis Cache INVALIDATE: user:%d, users:all (User updated)", user.ID)

	return nil
}

func (r *cachedUserRepository) Delete(ctx context.Context, id int64) error {
	err := r.baseRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	r.redisClient.Del(ctx, getCacheKey(id))
	r.redisClient.Del(ctx, "users:all")
	log.Printf("⚠ Redis Cache INVALIDATE: user:%d, users:all (User deleted)", id)

	return nil
}

func getCacheKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}
