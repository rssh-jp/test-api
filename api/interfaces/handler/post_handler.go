package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rssh-jp/test-api/api/usecase"
)

type PostHandler struct {
	postUsecase       usecase.PostUsecase // キャッシュ層を使う（デフォルト）
	directPostUsecase usecase.PostUsecase // キャッシュをバイパスしてDB直接アクセス
}

// NewPostHandler creates a new post handler
// postUsecase: キャッシュ層を経由するユースケース
// directPostUsecase: キャッシュをバイパスしてDB直接アクセスするユースケース
func NewPostHandler(postUsecase usecase.PostUsecase, directPostUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase:       postUsecase,
		directPostUsecase: directPostUsecase,
	}
}

// selectUsecase はクエリパラメータに応じて使用するユースケースを選択します
// no_cache=true が指定された場合はキャッシュをバイパス
func (h *PostHandler) selectUsecase(c echo.Context) usecase.PostUsecase {
	noCache := c.QueryParam("no_cache")
	if noCache == "true" || noCache == "1" {
		return h.directPostUsecase
	}
	return h.postUsecase
}

// GetPosts handles GET /posts
func (h *PostHandler) GetPosts(c echo.Context) error {
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	uc := h.selectUsecase(c) // クエリパラメータに応じてユースケースを選択

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	posts, total, err := uc.GetPosts(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts":    posts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetPostByID handles GET /posts/:id
func (h *PostHandler) GetPostByID(c echo.Context) error {
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	uc := h.selectUsecase(c) // クエリパラメータに応じてユースケースを選択

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid post ID",
		})
	}

	post, err := uc.GetPostByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Post not found",
		})
	}

	return c.JSON(http.StatusOK, post)
}

// GetPostBySlug handles GET /posts/slug/:slug
func (h *PostHandler) GetPostBySlug(c echo.Context) error {
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	uc := h.selectUsecase(c) // クエリパラメータに応じてユースケースを選択

	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Slug is required",
		})
	}

	post, err := uc.GetPostBySlug(ctx, slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Post not found",
		})
	}

	return c.JSON(http.StatusOK, post)
}

// GetPostsByCategory handles GET /posts/category/:slug
func (h *PostHandler) GetPostsByCategory(c echo.Context) error {
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	uc := h.selectUsecase(c) // クエリパラメータに応じてユースケースを選択

	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Category slug is required",
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	posts, err := uc.GetPostsByCategory(ctx, slug, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(http.StatusOK, posts)
}

// GetPostsByTag handles GET /posts/tag/:slug
func (h *PostHandler) GetPostsByTag(c echo.Context) error {
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	uc := h.selectUsecase(c) // クエリパラメータに応じてユースケースを選択

	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tag slug is required",
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}

	posts, err := uc.GetPostsByTag(ctx, slug, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(http.StatusOK, posts)
}

// GetFeaturedPosts handles GET /posts/featured
func (h *PostHandler) GetFeaturedPosts(c echo.Context) error {
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)
	uc := h.selectUsecase(c) // クエリパラメータに応じてユースケースを選択

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	posts, err := uc.GetFeaturedPosts(ctx, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve featured posts",
		})
	}

	return c.JSON(http.StatusOK, posts)
}
