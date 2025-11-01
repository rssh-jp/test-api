package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/rssh-jp/test-api/api/gen"
	redisCache "github.com/rssh-jp/test-api/api/infrastructure/cache/redis"
	mysqlRepo "github.com/rssh-jp/test-api/api/infrastructure/persistence/mysql"
	"github.com/rssh-jp/test-api/api/interfaces/handler"
	"github.com/rssh-jp/test-api/api/usecase"
)

func main() {
	// Load environment variables
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "mysql")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "testdb")

	redisHost := getEnv("REDIS_HOST", "redis")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	newrelicAppName := getEnv("NEW_RELIC_APP_NAME", "test-api")
	newrelicLicense := getEnv("NEW_RELIC_LICENSE_KEY", "")

	port := getEnv("PORT", "8080")

	// Initialize New Relic
	var nrApp *newrelic.Application
	var err error
	if newrelicLicense != "" {
		nrApp, err = newrelic.NewApplication(
			newrelic.ConfigAppName(newrelicAppName),
			newrelic.ConfigLicense(newrelicLicense),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			log.Printf("Warning: Failed to initialize New Relic: %v", err)
		} else {
			log.Println("New Relic initialized successfully")
		}
	} else {
		log.Println("Warning: NEW_RELIC_LICENSE_KEY not set, New Relic disabled")
	}

	// Initialize MySQL
	// Note: NewRelic instrumentation happens automatically via context
	// when nrecho middleware adds transaction to context
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	if nrApp != nil {
		log.Println("MySQL will be monitored by New Relic via context (from nrecho middleware)")
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Wait for database to be ready
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database connection... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to MySQL successfully")

	// Initialize Redis
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	defer redisClient.Close()

	// Add New Relic hook to Redis client
	if nrApp != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
		log.Println("Redis will be monitored by New Relic")
	}

	// Test Redis connection
	ctx := redisClient.Context()
	for i := 0; i < 30; i++ {
		_, err = redisClient.Ping(ctx).Result()
		if err == nil {
			break
		}
		log.Printf("Waiting for Redis connection... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")

	// Initialize repositories and services
	baseUserRepo := mysqlRepo.NewUserRepository(db)
	cachedUserRepo := redisCache.NewCachedUserRepository(baseUserRepo, redisClient)
	// ユーザーハンドラーにキャッシュ層とDB直接アクセス層の両方を渡す
	userUsecase := usecase.NewUserUsecase(cachedUserRepo)
	directUserUsecase := usecase.NewUserUsecase(baseUserRepo)
	
	// V2: フレームワーク非依存ハンドラーを作成し、ブリッジ経由でEchoに接続
	userHandlerV2 := handler.NewUserHandlerV2(userUsecase, directUserUsecase)
	userHandler := handler.NewUserHandlerBridge(userHandlerV2)

	// Initialize post-related services (complex JOIN queries with Redis cache)
	basePostRepo := mysqlRepo.NewPostRepository(db)
	cachedPostRepo := redisCache.NewCachedPostRepository(basePostRepo, redisClient)
	// 投稿ハンドラーにキャッシュ層とDB直接アクセス層の両方を渡す
	postUsecase := usecase.NewPostUsecase(cachedPostRepo)
	directPostUsecase := usecase.NewPostUsecase(basePostRepo)
	
	// V2: フレームワーク非依存ハンドラーを作成し、ブリッジ経由でEchoに接続
	postHandlerV2 := handler.NewPostHandlerV2(postUsecase, directPostUsecase)
	postHandler := handler.NewPostHandlerBridge(postHandlerV2)

	// Initialize user detail service (complex JOIN queries for all user-related data)
	userDetailRepo := mysqlRepo.NewUserDetailRepository(db)
	userDetailUsecase := usecase.NewUserDetailUsecase(userDetailRepo)
	
	// V2: フレームワーク非依存ハンドラーを作成し、ブリッジ経由でEchoに接続
	userDetailHandlerV2 := handler.NewUserDetailHandlerV2(userDetailUsecase)
	userDetailHandler := handler.NewUserDetailHandlerBridge(userDetailHandlerV2)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// New Relic middleware
	if nrApp != nil {
		e.Use(nrecho.Middleware(nrApp))
	}

	// Register routes using OpenAPI generated code
	gen.RegisterHandlers(e, userHandler)

	// Register post routes (complex JOIN queries)
	e.GET("/posts", postHandler.GetPosts)
	e.GET("/posts/featured", postHandler.GetFeaturedPosts)
	e.GET("/posts/:id", postHandler.GetPostByID)
	e.GET("/posts/slug/:slug", postHandler.GetPostBySlug)
	e.GET("/posts/category/:slug", postHandler.GetPostsByCategory)
	e.GET("/posts/tag/:slug", postHandler.GetPostsByTag)

	// Register user detail routes (complex JOIN queries for all user-related data)
	e.GET("/users/:id/detail", userDetailHandler.GetUserDetailByID)
	e.GET("/users/username/:username/detail", userDetailHandler.GetUserDetailByUsername)

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
