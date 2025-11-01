package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

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
