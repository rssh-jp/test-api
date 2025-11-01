package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rssh-jp/test-api/api/gen"
)

// UserHandlerBridge はEchoのServerInterfaceとフレームワーク非依存ハンドラーを繋ぐブリッジ
// このブリッジを通して、Echo固有の処理をフレームワーク非依存のハンドラーに委譲します
// Echo Context → HTTPContext の変換もこの層で行います
type UserHandlerBridge struct {
	handler *UserHandlerV2
}

// NewUserHandlerBridge creates a new bridge that implements gen.ServerInterface
func NewUserHandlerBridge(handler *UserHandlerV2) gen.ServerInterface {
	return &UserHandlerBridge{
		handler: handler,
	}
}

// echoHTTPContext はEcho Contextをフレームワーク非依存のHTTPContextに変換する
// Bridge内でのみ使用される内部実装
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
