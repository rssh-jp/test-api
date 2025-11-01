package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

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

// parseNoCache はクエリパラメータからno_cacheを解析します
func parseNoCache(c echo.Context) bool {
	noCache := c.QueryParam("no_cache")
	return noCache == "true" || noCache == "1"
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
