package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rssh-jp/test-api/api/gen"
	"github.com/rssh-jp/test-api/api/usecase"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) gen.ServerInterface {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// HealthCheck implements the health check endpoint
func (h *UserHandler) HealthCheck(ctx echo.Context) error {
	response := gen.HealthResponse{
		Status:  "healthy",
		Message: "Service is running",
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetUsers implements get all users endpoint
func (h *UserHandler) GetUsers(ctx echo.Context) error {
	txn := newrelic.FromContext(ctx.Request().Context())
	reqCtx := newrelic.NewContext(ctx.Request().Context(), txn)
	users, err := h.userUsecase.GetAllUsers(reqCtx)
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

// GetUserById implements get user by ID endpoint
func (h *UserHandler) GetUserById(ctx echo.Context, id int64) error {
	txn := newrelic.FromContext(ctx.Request().Context())
	reqCtx := newrelic.NewContext(ctx.Request().Context(), txn)
	user, err := h.userUsecase.GetUserByID(reqCtx, id)
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

// CreateUser implements create user endpoint
func (h *UserHandler) CreateUser(ctx echo.Context) error {
	var req gen.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
		})
	}

	txn := newrelic.FromContext(ctx.Request().Context())
	reqCtx := newrelic.NewContext(ctx.Request().Context(), txn)
	user, err := h.userUsecase.CreateUser(reqCtx, req.Name, string(req.Email), req.Age)
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

// UpdateUser implements update user endpoint
func (h *UserHandler) UpdateUser(ctx echo.Context, id int64) error {
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

	txn := newrelic.FromContext(ctx.Request().Context())
	reqCtx := newrelic.NewContext(ctx.Request().Context(), txn)

	user, err := h.userUsecase.UpdateUser(reqCtx, id, req.Name, email, req.Age)
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

// DeleteUser implements delete user endpoint
func (h *UserHandler) DeleteUser(ctx echo.Context, id int64) error {
	txn := newrelic.FromContext(ctx.Request().Context())
	reqCtx := newrelic.NewContext(ctx.Request().Context(), txn)

	err := h.userUsecase.DeleteUser(reqCtx, id)
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

// PathToID is a helper function to convert path parameter to int64
func PathToID(ctx echo.Context, paramName string) (int64, error) {
	idStr := ctx.Param(paramName)
	return strconv.ParseInt(idStr, 10, 64)
}
