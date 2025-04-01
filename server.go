package echox

import (
	"errors"
	"net/http"
	"time"

	"github.com/boostgo/appx"
	"github.com/boostgo/errorx"
	"github.com/boostgo/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type route struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

var (
	_routes = make([]route, 0)
)

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

func run(address string) error {
	handler := echo.New()

	handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Auth-Token"},
		AllowCredentials: true,
	}))
	handler.Use(RecoverMiddleware())

	handler.RouteNotFound("*", func(ctx echo.Context) error {
		return Error(ctx, errorx.
			New("Route not found").
			SetError(errorx.ErrNotFound).
			AddContext("url", ctx.Request().RequestURI))
	})

	if getTracer().AmIMaster() {
		handler.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: uuid.NewString,
			RequestIDHandler: func(ctx echo.Context, traceID string) {
				ctx.SetRequest(
					ctx.Request().WithContext(getTracer().Set(ctx.Request().Context())))
			},
			TargetHeader: "X-Trace-ID",
		}))
	}

	// set middlewares
	for _, mid := range _middlewares {
		handler.Use(mid)
	}

	// set routes
	for _, r := range _routes {
		handler.Add(r.Method, r.Path, r.Handler, r.Middlewares...)
	}

	appx.Tear(func() error {
		return handler.Shutdown(appx.Context())
	})

	if err := handler.Start(address); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return errorx.
			New("Start server").
			SetError(err)
	}

	return nil
}

func Run(address string, waitTime ...time.Duration) {
	go func() {
		if err := run(address); err != nil {
			log.
				Error().
				Err(err).
				Msg("Run server")
			appx.Cancel()
		}
	}()

	appx.GracefulLog(func() {
		log.
			Info().
			Msg("Graceful shutdown...")
	})
	appx.Wait(waitTime...)
}
