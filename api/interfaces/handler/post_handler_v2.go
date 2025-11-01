package handler

import (
	"net/http"

	"github.com/rssh-jp/test-api/api/usecase"
)

// PostHandlerV2 はフレームワーク非依存の投稿ハンドラー
type PostHandlerV2 struct {
	postUsecase       usecase.PostUsecase // キャッシュ層を使う（デフォルト）
	directPostUsecase usecase.PostUsecase // キャッシュをバイパスしてDB直接アクセス
}

// NewPostHandlerV2 creates a new framework-independent post handler
func NewPostHandlerV2(postUsecase usecase.PostUsecase, directPostUsecase usecase.PostUsecase) *PostHandlerV2 {
	return &PostHandlerV2{
		postUsecase:       postUsecase,
		directPostUsecase: directPostUsecase,
	}
}

// selectUsecase はクエリパラメータに応じて使用するユースケースを選択します
func (h *PostHandlerV2) selectUsecase(noCache bool) usecase.PostUsecase {
	if noCache {
		return h.directPostUsecase
	}
	return h.postUsecase
}

// GetPosts は投稿一覧を取得します（フレームワーク非依存）
func (h *PostHandlerV2) GetPosts(ctx HTTPContext, page, pageSize int, noCache bool) error {
	reqCtx := ctx.Context()
	uc := h.selectUsecase(noCache)

	posts, total, err := uc.GetPosts(reqCtx, page, pageSize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve posts",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"posts":    posts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetPostByID はIDで投稿を取得します（フレームワーク非依存）
func (h *PostHandlerV2) GetPostByID(ctx HTTPContext, id int64, noCache bool) error {
	reqCtx := ctx.Context()
	uc := h.selectUsecase(noCache)

	post, err := uc.GetPostByID(reqCtx, id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Post not found",
		})
	}

	return ctx.JSON(http.StatusOK, post)
}

// GetPostBySlug はスラッグで投稿を取得します（フレームワーク非依存）
func (h *PostHandlerV2) GetPostBySlug(ctx HTTPContext, slug string, noCache bool) error {
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Slug is required",
		})
	}

	reqCtx := ctx.Context()
	uc := h.selectUsecase(noCache)

	post, err := uc.GetPostBySlug(reqCtx, slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Post not found",
		})
	}

	return ctx.JSON(http.StatusOK, post)
}

// GetPostsByCategory はカテゴリー別に投稿を取得します（フレームワーク非依存）
func (h *PostHandlerV2) GetPostsByCategory(ctx HTTPContext, slug string, page, pageSize int, noCache bool) error {
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Category slug is required",
		})
	}

	reqCtx := ctx.Context()
	uc := h.selectUsecase(noCache)

	posts, err := uc.GetPostsByCategory(reqCtx, slug, page, pageSize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve posts",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}

// GetPostsByTag はタグ別に投稿を取得します（フレームワーク非依存）
func (h *PostHandlerV2) GetPostsByTag(ctx HTTPContext, slug string, page, pageSize int, noCache bool) error {
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tag slug is required",
		})
	}

	reqCtx := ctx.Context()
	uc := h.selectUsecase(noCache)

	posts, err := uc.GetPostsByTag(reqCtx, slug, page, pageSize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve posts",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}

// GetFeaturedPosts は注目投稿を取得します（フレームワーク非依存）
func (h *PostHandlerV2) GetFeaturedPosts(ctx HTTPContext, limit int, noCache bool) error {
	reqCtx := ctx.Context()
	uc := h.selectUsecase(noCache)

	posts, err := uc.GetFeaturedPosts(reqCtx, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve featured posts",
		})
	}

	return ctx.JSON(http.StatusOK, posts)
}
