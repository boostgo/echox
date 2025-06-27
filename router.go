package echox

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

var _routes = make([]route, 0)

type Router interface {
	Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) []*echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Use(m ...echo.MiddlewareFunc)
	RouteNotFound(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
}

func RegisterRoute(method, path string, handlerFunc echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	_routes = append(_routes, route{
		Method:      method,
		Path:        path,
		Handler:     handlerFunc,
		Middlewares: m,
	})
}

func GET(path string, handlerFunc echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	RegisterRoute(http.MethodGet, path, handlerFunc, m...)
}

func POST(path string, handlerFunc echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	RegisterRoute(http.MethodPost, path, handlerFunc, m...)
}

func PUT(path string, handlerFunc echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	RegisterRoute(http.MethodPut, path, handlerFunc, m...)
}

func PATCH(path string, handlerFunc echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	RegisterRoute(http.MethodPatch, path, handlerFunc, m...)
}

func DELETE(path string, handlerFunc echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	RegisterRoute(http.MethodDelete, path, handlerFunc, m...)
}
