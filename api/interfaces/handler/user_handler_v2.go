package handler

import (
	"database/sql"
	"errors"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rssh-jp/test-api/api/gen"
	"github.com/rssh-jp/test-api/api/usecase"
)

// UserHandlerV2 はフレームワーク非依存のユーザーハンドラー
// HTTPContextインターフェースを使用することで、Echo/Gin/Chiなど
// 任意のフレームワークに対応できます
type UserHandlerV2 struct {
	userUsecase       usecase.UserUsecase // キャッシュ層を使う（デフォルト）
	directUserUsecase usecase.UserUsecase // キャッシュをバイパスしてDB直接アクセス
}

// NewUserHandlerV2 creates a new framework-independent user handler
func NewUserHandlerV2(userUsecase usecase.UserUsecase, directUserUsecase usecase.UserUsecase) *UserHandlerV2 {
	return &UserHandlerV2{
		userUsecase:       userUsecase,
		directUserUsecase: directUserUsecase,
	}
}

// GetUsers は全ユーザーを取得します（フレームワーク非依存）
func (h *UserHandlerV2) GetUsers(ctx HTTPContext, params gen.GetUsersParams) error {
	reqCtx := ctx.Context()
	
	// 型安全なパラメータから判定
	uc := h.userUsecase
	if params.NoCache != nil && *params.NoCache {
		uc = h.directUserUsecase
	}

	users, err := uc.GetAllUsers(reqCtx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to retrieve users",
		})
	}

	// Convert domain users to API users
	apiUsers := make([]gen.User, len(users))
	for i, user := range users {
		apiUsers[i] = gen.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     openapi_types.Email(user.Email),
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, apiUsers)
}

// GetUserById はIDでユーザーを取得します（フレームワーク非依存）
func (h *UserHandlerV2) GetUserById(ctx HTTPContext, id int64, params gen.GetUserByIdParams) error {
	reqCtx := ctx.Context()
	
	// 型安全なパラメータから判定
	uc := h.userUsecase
	if params.NoCache != nil && *params.NoCache {
		uc = h.directUserUsecase
	}

	user, err := uc.GetUserByID(reqCtx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, gen.Error{
				Message: "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to retrieve user",
		})
	}

	apiUser := gen.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, apiUser)
}

// CreateUser はユーザーを作成します（フレームワーク非依存）
func (h *UserHandlerV2) CreateUser(ctx HTTPContext) error {
	var req gen.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
		})
	}

	reqCtx := ctx.Context()
	// 作成時は常にキャッシュ層を使用（書き込み操作）
	uc := h.userUsecase

	user, err := uc.CreateUser(reqCtx, req.Name, string(req.Email), req.Age)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to create user",
		})
	}

	apiUser := gen.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return ctx.JSON(http.StatusCreated, apiUser)
}

// UpdateUser はユーザーを更新します（フレームワーク非依存）
func (h *UserHandlerV2) UpdateUser(ctx HTTPContext, id int64) error {
	var req gen.UpdateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
		})
	}

	var email *string
	if req.Email != nil {
		emailStr := string(*req.Email)
		email = &emailStr
	}

	reqCtx := ctx.Context()
	// 更新時は常にキャッシュ層を使用（書き込み操作）
	uc := h.userUsecase

	user, err := uc.UpdateUser(reqCtx, id, req.Name, email, req.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, gen.Error{
				Message: "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to update user",
		})
	}

	apiUser := gen.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     openapi_types.Email(user.Email),
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, apiUser)
}

// DeleteUser はユーザーを削除します（フレームワーク非依存）
func (h *UserHandlerV2) DeleteUser(ctx HTTPContext, id int64) error {
	reqCtx := ctx.Context()
	// 削除時は常にキャッシュ層を使用（書き込み操作）
	uc := h.userUsecase

	err := uc.DeleteUser(reqCtx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, gen.Error{
				Message: "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to delete user",
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// HealthCheck はヘルスチェックを実行します（フレームワーク非依存）
func (h *UserHandlerV2) HealthCheck(ctx HTTPContext) error {
	response := gen.HealthResponse{
		Status:  "healthy",
		Message: "Service is running",
	}
	return ctx.JSON(http.StatusOK, response)
}
