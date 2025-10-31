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

type cachedPostRepository struct {
	baseRepo    domain.PostRepository
	redisClient *redis.Client
	ctx         context.Context
	ttl         time.Duration
}

// NewCachedPostRepository creates a new cached post repository
func NewCachedPostRepository(baseRepo domain.PostRepository, redisClient *redis.Client) domain.PostRepository {
	return &cachedPostRepository{
		baseRepo:    baseRepo,
		redisClient: redisClient,
		ctx:         context.Background(),
		ttl:         5 * time.Minute, // Cache TTL: 5 minutes
	}
}

func (r *cachedPostRepository) FindAllWithDetails(ctx context.Context, limit, offset int) ([]domain.PostWithDetails, error) {
	cacheKey := fmt.Sprintf("posts:all:limit=%d:offset=%d", limit, offset)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []domain.PostWithDetails
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			log.Printf("✓ Redis Cache HIT: %s (4-table JOIN cached)", cacheKey)
			return posts, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL (4-table JOIN)", cacheKey)
	posts, err := r.baseRepo.FindAllWithDetails(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(posts)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return posts, nil
}

func (r *cachedPostRepository) FindByIDWithDetails(ctx context.Context, id int64) (*domain.PostWithDetails, error) {
	cacheKey := getPostCacheKey(id)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var post domain.PostWithDetails
		if err := json.Unmarshal([]byte(cached), &post); err == nil {
			log.Printf("✓ Redis Cache HIT: %s (multi-table JOIN)", cacheKey)
			return &post, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL (multi-table JOIN)", cacheKey)
	post, err := r.baseRepo.FindByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(post)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return post, nil
}

func (r *cachedPostRepository) FindBySlugWithDetails(ctx context.Context, slug string) (*domain.PostWithDetails, error) {
	cacheKey := fmt.Sprintf("post:slug:%s", slug)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var post domain.PostWithDetails
		if err := json.Unmarshal([]byte(cached), &post); err == nil {
			log.Printf("✓ Redis Cache HIT: %s (multi-table JOIN)", cacheKey)
			return &post, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL (multi-table JOIN)", cacheKey)
	post, err := r.baseRepo.FindBySlugWithDetails(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(post)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return post, nil
}

func (r *cachedPostRepository) FindByCategoryWithDetails(ctx context.Context, categorySlug string, limit, offset int) ([]domain.PostWithDetails, error) {
	cacheKey := fmt.Sprintf("posts:category:%s:limit=%d:offset=%d", categorySlug, limit, offset)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []domain.PostWithDetails
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			log.Printf("✓ Redis Cache HIT: %s (category JOIN cached)", cacheKey)
			return posts, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL (category JOIN)", cacheKey)
	posts, err := r.baseRepo.FindByCategoryWithDetails(ctx, categorySlug, limit, offset)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(posts)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return posts, nil
}

func (r *cachedPostRepository) FindByTagWithDetails(ctx context.Context, tagSlug string, limit, offset int) ([]domain.PostWithDetails, error) {
	cacheKey := fmt.Sprintf("posts:tag:%s:limit=%d:offset=%d", tagSlug, limit, offset)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []domain.PostWithDetails
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			log.Printf("✓ Redis Cache HIT: %s (tag JOIN cached)", cacheKey)
			return posts, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL (tag JOIN)", cacheKey)
	posts, err := r.baseRepo.FindByTagWithDetails(ctx, tagSlug, limit, offset)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(posts)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return posts, nil
}

func (r *cachedPostRepository) FindFeaturedWithDetails(ctx context.Context, limit int) ([]domain.PostWithDetails, error) {
	cacheKey := fmt.Sprintf("posts:featured:limit=%d", limit)

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []domain.PostWithDetails
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			log.Printf("✓ Redis Cache HIT: %s (featured posts cached)", cacheKey)
			return posts, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL (featured posts)", cacheKey)
	posts, err := r.baseRepo.FindFeaturedWithDetails(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Store in cache
	data, _ := json.Marshal(posts)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return posts, nil
}

func (r *cachedPostRepository) GetTotalCount(ctx context.Context) (int64, error) {
	cacheKey := "posts:totalcount"

	// Try to get from cache
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var count int64
		if err := json.Unmarshal([]byte(cached), &count); err == nil {
			log.Printf("✓ Redis Cache HIT: %s", cacheKey)
			return count, nil
		}
	}

	// Cache miss, get from database
	log.Printf("✗ Redis Cache MISS: %s - Fetching from MySQL", cacheKey)
	count, err := r.baseRepo.GetTotalCount(ctx)
	if err != nil {
		return 0, err
	}

	// Store in cache
	data, _ := json.Marshal(count)
	r.redisClient.Set(ctx, cacheKey, data, r.ttl)
	log.Printf("→ Redis Cache SET: %s (TTL: %v)", cacheKey, r.ttl)

	return count, nil
}

func (r *cachedPostRepository) IncrementViewCount(ctx context.Context, postID int64) error {
	// Increment view count in database
	err := r.baseRepo.IncrementViewCount(ctx, postID)
	if err != nil {
		return err
	}

	// Invalidate related caches
	r.redisClient.Del(ctx, getPostCacheKey(postID))
	// Invalidate list caches
	keys, _ := r.redisClient.Keys(ctx, "posts:*").Result()
	if len(keys) > 0 {
		r.redisClient.Del(ctx, keys...)
	}
	log.Printf("⚠ Redis Cache INVALIDATE: post:%d (view count incremented)", postID)

	return nil
}

func getPostCacheKey(id int64) string {
	return fmt.Sprintf("post:%d", id)
}
