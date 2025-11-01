package handler

import (
	"net/http"

	"github.com/rssh-jp/test-api/api/usecase"
)

// UserDetailHandlerV2 はフレームワーク非依存のユーザー詳細ハンドラー
type UserDetailHandlerV2 struct {
	usecase usecase.UserDetailUsecase
}

// NewUserDetailHandlerV2 creates a new framework-independent user detail handler
func NewUserDetailHandlerV2(usecase usecase.UserDetailUsecase) *UserDetailHandlerV2 {
	return &UserDetailHandlerV2{usecase: usecase}
}

// GetUserDetailByID は指定されたIDのユーザーの全関連情報を取得します（フレームワーク非依存）
func (h *UserDetailHandlerV2) GetUserDetailByID(ctx HTTPContext, id int64) error {
	reqCtx := ctx.Context()

	// ユーザー詳細情報を取得
	detail, err := h.usecase.GetUserDetailByID(reqCtx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user details",
		})
	}

	return ctx.JSON(http.StatusOK, detail)
}

// GetUserDetailByUsername は指定されたユーザー名のユーザーの全関連情報を取得します（フレームワーク非依存）
func (h *UserDetailHandlerV2) GetUserDetailByUsername(ctx HTTPContext, username string) error {
	if username == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Username is required",
		})
	}

	reqCtx := ctx.Context()

	// ユーザー詳細情報を取得
	detail, err := h.usecase.GetUserDetailByUsername(reqCtx, username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user details",
		})
	}

	return ctx.JSON(http.StatusOK, detail)
}
