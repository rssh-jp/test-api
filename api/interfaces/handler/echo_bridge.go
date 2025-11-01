package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rssh-jp/test-api/api/gen"
)

// ============================================================================
// Echo HTTPContext Adapter (共通実装)
// ============================================================================

// echoHTTPContext はEcho Contextをフレームワーク非依存のHTTPContextに変換する
type echoHTTPContext struct {
	ctx echo.Context
}

// newEchoHTTPContext creates a new Echo HTTP context adapter
func newEchoHTTPContext(ctx echo.Context) HTTPContext {
	return &echoHTTPContext{ctx: ctx}
}

// Context returns the request context with NewRelic transaction
func (e *echoHTTPContext) Context() context.Context {
	// NewRelicトランザクションをコンテキストに注入
	txn := newrelic.FromContext(e.ctx.Request().Context())
	return newrelic.NewContext(e.ctx.Request().Context(), txn)
}

// Request returns the underlying *http.Request
func (e *echoHTTPContext) Request() *http.Request {
	return e.ctx.Request()
}

// Response returns the underlying http.ResponseWriter
func (e *echoHTTPContext) Response() http.ResponseWriter {
	return e.ctx.Response().Writer
}

// Bind binds the request body to the given struct
func (e *echoHTTPContext) Bind(i interface{}) error {
	return e.ctx.Bind(i)
}

// JSON sends a JSON response with the given status code
func (e *echoHTTPContext) JSON(code int, data interface{}) error {
	return e.ctx.JSON(code, data)
}

// NoContent sends a response with no body
func (e *echoHTTPContext) NoContent(code int) error {
	return e.ctx.NoContent(code)
}

// QueryParam returns the query parameter value by name
func (e *echoHTTPContext) QueryParam(name string) string {
	return e.ctx.QueryParam(name)
}

// Param returns the path parameter value by name
func (e *echoHTTPContext) Param(name string) string {
	return e.ctx.Param(name)
}

// parseNoCache はクエリパラメータからno_cacheを解析します（共通ヘルパー）
func parseNoCache(c echo.Context) bool {
	noCache := c.QueryParam("no_cache")
	return noCache == "true" || noCache == "1"
}

// ============================================================================
// UserHandlerBridge (OpenAPI生成インターフェース実装)
// ============================================================================

// UserHandlerBridge はEchoのServerInterfaceとフレームワーク非依存ハンドラーを繋ぐブリッジ
type UserHandlerBridge struct {
	handler *UserHandlerV2
}

// NewUserHandlerBridge creates a new bridge that implements gen.ServerInterface
func NewUserHandlerBridge(handler *UserHandlerV2) gen.ServerInterface {
	return &UserHandlerBridge{
		handler: handler,
	}
}

// HealthCheck implements the health check endpoint (Echo → Framework-independent)
func (b *UserHandlerBridge) HealthCheck(ctx echo.Context) error {
	httpCtx := newEchoHTTPContext(ctx)
	return b.handler.HealthCheck(httpCtx)
}

// GetUsers implements get all users endpoint (Echo → Framework-independent)
func (b *UserHandlerBridge) GetUsers(ctx echo.Context, params gen.GetUsersParams) error {
	httpCtx := newEchoHTTPContext(ctx)
	return b.handler.GetUsers(httpCtx, params)
}

// GetUserById implements get user by ID endpoint (Echo → Framework-independent)
func (b *UserHandlerBridge) GetUserById(ctx echo.Context, id int64, params gen.GetUserByIdParams) error {
	httpCtx := newEchoHTTPContext(ctx)
	return b.handler.GetUserById(httpCtx, id, params)
}

// CreateUser implements create user endpoint (Echo → Framework-independent)
func (b *UserHandlerBridge) CreateUser(ctx echo.Context) error {
	httpCtx := newEchoHTTPContext(ctx)
	return b.handler.CreateUser(httpCtx)
}

// UpdateUser implements update user endpoint (Echo → Framework-independent)
func (b *UserHandlerBridge) UpdateUser(ctx echo.Context, id int64) error {
	httpCtx := newEchoHTTPContext(ctx)
	return b.handler.UpdateUser(httpCtx, id)
}

// DeleteUser implements delete user endpoint (Echo → Framework-independent)
func (b *UserHandlerBridge) DeleteUser(ctx echo.Context, id int64) error {
	httpCtx := newEchoHTTPContext(ctx)
	return b.handler.DeleteUser(httpCtx, id)
}

// ============================================================================
// UserDetailHandlerBridge
// ============================================================================

// UserDetailHandlerBridge はEchoとフレームワーク非依存UserDetailHandlerを繋ぐブリッジ
type UserDetailHandlerBridge struct {
	handler *UserDetailHandlerV2
}

// NewUserDetailHandlerBridge creates a new bridge for user detail handler
func NewUserDetailHandlerBridge(handler *UserDetailHandlerV2) *UserDetailHandlerBridge {
	return &UserDetailHandlerBridge{
		handler: handler,
	}
}

// GetUserDetailByID は指定されたIDのユーザーの全関連情報を取得します (Echo → Framework-independent)
// GET /users/:id/detail
func (b *UserDetailHandlerBridge) GetUserDetailByID(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	// パスパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(400, map[string]string{
			"error": "Invalid user ID",
		})
	}

	return b.handler.GetUserDetailByID(httpCtx, id)
}

// GetUserDetailByUsername は指定されたユーザー名のユーザーの全関連情報を取得します (Echo → Framework-independent)
// GET /users/username/:username/detail
func (b *UserDetailHandlerBridge) GetUserDetailByUsername(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	// パスパラメータからユーザー名を取得
	username := c.Param("username")

	return b.handler.GetUserDetailByUsername(httpCtx, username)
}

// ============================================================================
// PostHandlerBridge
// ============================================================================

// PostHandlerBridge はEchoとフレームワーク非依存PostHandlerを繋ぐブリッジ
type PostHandlerBridge struct {
	handler *PostHandlerV2
}

// NewPostHandlerBridge creates a new bridge for post handler
func NewPostHandlerBridge(handler *PostHandlerV2) *PostHandlerBridge {
	return &PostHandlerBridge{
		handler: handler,
	}
}

// GetPosts handles GET /posts (Echo → Framework-independent)
func (b *PostHandlerBridge) GetPosts(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	noCache := parseNoCache(c)

	return b.handler.GetPosts(httpCtx, page, pageSize, noCache)
}

// GetPostByID handles GET /posts/:id (Echo → Framework-independent)
func (b *PostHandlerBridge) GetPostByID(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(400, map[string]string{
			"error": "Invalid post ID",
		})
	}

	noCache := parseNoCache(c)

	return b.handler.GetPostByID(httpCtx, id, noCache)
}

// GetPostBySlug handles GET /posts/slug/:slug (Echo → Framework-independent)
func (b *PostHandlerBridge) GetPostBySlug(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	slug := c.Param("slug")
	noCache := parseNoCache(c)

	return b.handler.GetPostBySlug(httpCtx, slug, noCache)
}

// GetPostsByCategory handles GET /posts/category/:slug (Echo → Framework-independent)
func (b *PostHandlerBridge) GetPostsByCategory(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	slug := c.Param("slug")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	noCache := parseNoCache(c)

	return b.handler.GetPostsByCategory(httpCtx, slug, page, pageSize, noCache)
}

// GetPostsByTag handles GET /posts/tag/:slug (Echo → Framework-independent)
func (b *PostHandlerBridge) GetPostsByTag(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	slug := c.Param("slug")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	noCache := parseNoCache(c)

	return b.handler.GetPostsByTag(httpCtx, slug, page, pageSize, noCache)
}

// GetFeaturedPosts handles GET /posts/featured (Echo → Framework-independent)
func (b *PostHandlerBridge) GetFeaturedPosts(c echo.Context) error {
	httpCtx := newEchoHTTPContext(c)

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	noCache := parseNoCache(c)

	return b.handler.GetFeaturedPosts(httpCtx, limit, noCache)
}
