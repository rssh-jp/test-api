package handler

import (
	"context"
	"net/http"
)

// HTTPContext はフレームワーク非依存のHTTPコンテキストインターフェース
// このインターフェースを使うことで、Echo/Gin/Chiなどのフレームワークを
// 容易に切り替えることができます
type HTTPContext interface {
	// Context returns the standard context.Context
	Context() context.Context

	// Request returns the underlying *http.Request
	Request() *http.Request

	// Response returns the underlying http.ResponseWriter
	Response() http.ResponseWriter

	// Bind binds the request body to the given struct
	Bind(interface{}) error

	// JSON sends a JSON response with the given status code
	JSON(code int, data interface{}) error

	// NoContent sends a response with no body
	NoContent(code int) error

	// QueryParam returns the query parameter value by name
	QueryParam(name string) string

	// Param returns the path parameter value by name
	Param(name string) string
}
