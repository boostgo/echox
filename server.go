package echox

import (
	"errors"
	"net/http"
	"time"

	"github.com/boostgo/appx"
	"github.com/boostgo/configx"
	"github.com/boostgo/convert"
	"github.com/boostgo/httpx"
	"github.com/boostgo/log"
	"github.com/boostgo/trace"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func run(address string) error {
	handler := echo.New()

	// add CORS middleware
	handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Auth-Token"},
		AllowCredentials: true,
	}))

	// add recover middleware
	handler.Use(RecoverMiddleware())

	// register not found route
	handler.RouteNotFound("*", func(ctx echo.Context) error {
		return Error(ctx, newRouteNotFoundError(ctx.Request()))
	})

	// add trace middleware
	if trace.AmIMaster() {
		handler.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: uuid.NewString,
			RequestIDHandler: func(ctx echo.Context, traceID string) {
				requestCtx := ctx.Request().Context()
				requestCtx = trace.SetID(requestCtx, traceID)
				ctx.SetRequest(ctx.Request().WithContext(requestCtx))
			},
			TargetHeader: TraceKey,
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

	// set groups
	for _, g := range _groups {
		group := handler.Group(g.basePath, g.middlewares...)
		for _, r := range g.routes {
			group.Add(r.Method, r.Path, r.Handler, r.Middlewares...)
		}
	}

	// add server shutdown teardown func
	appx.Tear(func() error {
		return handler.Shutdown(appx.Context())
	})

	// print all registered routes (only in dev mode)
	if configx.Production() {
		log.
			Info().
			Int("routes_count", len(handler.Routes())).
			Msg("Registered routes")

		for idx, r := range handler.Routes() {
			log.
				Info().
				Str("method", r.Method).
				Str("path", r.Path).
				Str("name", r.Name).
				Msg(convert.StringFromInt(idx+1) + ". Registered route")
		}
	}

	// start server
	if err := handler.Start(address); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return httpx.ErrStartServer.SetError(err)
	}

	return nil
}

func Run(address string, waitTime ...time.Duration) {
	// run server in new goroutine
	go func() {
		if err := run(address); err != nil {
			log.
				Error().
				Err(err).
				Msg("Run server")

			// if server run failed - call app shutdown
			appx.Cancel()
		}
	}()

	// add app graceful shutdown log
	appx.GracefulLog(func() {
		log.
			Info().
			Msg("Graceful shutdown...")
	})

	// wait till the end of app lifetime
	appx.Wait(waitTime...)
}
