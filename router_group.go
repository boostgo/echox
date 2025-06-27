package echox

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	PROPFIND = "PROPFIND"
	REPORT   = "REPORT"
)

var (
	_groups  = make([]*RouterGroup, 0)
	_methods = []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		PROPFIND,
		http.MethodPut,
		http.MethodTrace,
		REPORT,
	}
)

type RouterGroup struct {
	basePath    string
	routes      []route
	middlewares []echo.MiddlewareFunc
}

func Group(basePath string, m ...echo.MiddlewareFunc) *RouterGroup {
	g := &RouterGroup{
		basePath:    basePath,
		routes:      make([]route, 0),
		middlewares: m,
	}
	_groups = append(_groups, g)
	return g
}

func (g *RouterGroup) Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	for _, method := range _methods {
		g.routes = append(g.routes, route{
			Method:      method,
			Path:        path,
			Handler:     h,
			Middlewares: m,
		})
	}
}

func (g *RouterGroup) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodPost,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodGet,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodPut,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodPatch,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodDelete,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodHead,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      http.MethodOptions,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) RouteNotFound(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      echo.RouteNotFound,
		Path:        path,
		Handler:     h,
		Middlewares: m,
	})
}

func (g *RouterGroup) Use(m ...echo.MiddlewareFunc) {
	g.middlewares = append(g.middlewares, m...)
}

func (g *RouterGroup) Group(prefix string, m ...echo.MiddlewareFunc) *RouterGroup {
	return Group(g.basePath+prefix, append(g.middlewares, m...)...)
}

func (g *RouterGroup) Register(method, url string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.routes = append(g.routes, route{
		Method:      method,
		Path:        url,
		Handler:     h,
		Middlewares: m,
	})
}
