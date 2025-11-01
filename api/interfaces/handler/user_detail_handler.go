package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rssh-jp/test-api/api/usecase"
)

type UserDetailHandler struct {
	usecase usecase.UserDetailUsecase
}

// NewUserDetailHandler creates a new user detail handler
func NewUserDetailHandler(usecase usecase.UserDetailUsecase) *UserDetailHandler {
	return &UserDetailHandler{usecase: usecase}
}

// GetUserDetailByID は指定されたIDのユーザーの全関連情報を取得します
// GET /users/:id/detail
func (h *UserDetailHandler) GetUserDetailByID(c echo.Context) error {
	// NewRelicトランザクションコンテキスト
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)

	// パスパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// ユーザー詳細情報を取得
	detail, err := h.usecase.GetUserDetailByID(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user details",
		})
	}

	return c.JSON(http.StatusOK, detail)
}

// GetUserDetailByUsername は指定されたユーザー名のユーザーの全関連情報を取得します
// GET /users/username/:username/detail
func (h *UserDetailHandler) GetUserDetailByUsername(c echo.Context) error {
	// NewRelicトランザクションコンテキスト
	txn := newrelic.FromContext(c.Request().Context())
	ctx := newrelic.NewContext(c.Request().Context(), txn)

	// パスパラメータからユーザー名を取得
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Username is required",
		})
	}

	// ユーザー詳細情報を取得
	detail, err := h.usecase.GetUserDetailByUsername(ctx, username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user details",
		})
	}

	return c.JSON(http.StatusOK, detail)
}
