package echox

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/boostgo/convert"
	"github.com/boostgo/errorx"
	"github.com/boostgo/httpx"
	"github.com/boostgo/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	rawResponseKey = "response-raw"
)

var (
	_middlewares = make([]echo.MiddlewareFunc, 0)
)

func RegisterMiddleware(mid echo.MiddlewareFunc) {
	_middlewares = append(_middlewares, mid)
}

func RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := errorx.Try(func() error {
				return next(ctx)
			}); err != nil {
				return Error(ctx, err)
			}

			return nil
		}
	}
}

func TimeoutMiddleware(duration time.Duration) echo.MiddlewareFunc {
	return middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		ErrorHandler: func(err error, ctx echo.Context) error {
			return Error(
				ctx,
				errorx.
					New("Request reached timeout").
					SetError(err, errorx.ErrTimeout),
			)
		},
		Timeout: duration,
	})
}

// Raw if middleware set, all responses by this middleware will be returned in "raw" way (no successOutput object)
func Raw() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			localCtx := Context(ctx)
			localCtx = context.WithValue(localCtx, rawResponseKey, true)
			ctx.SetRequest(ctx.Request().WithContext(localCtx))
			return next(ctx)
		}
	}
}

func Cache(ttl time.Duration, distributor httpx.CacheDistributor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// try load response from cache
			responseBody, cacheOk, err := distributor.Get(Context(ctx), ctx.Request())
			if err != nil {
				cacheOk = false

				if !errors.Is(err, errorx.ErrNotFound) {
					log.
						Error(Context(ctx)).
						Err(err).
						Msg("Get cache by HTTP distributor")
				}
			}

			// return cached response
			if cacheOk {
				return SuccessRaw(ctx, http.StatusOK, responseBody, httpx.ContentTypeJSON)
			}

			// call handler method to generate response
			response := ctx.Response()
			var responseBuffer bytes.Buffer
			mw := io.MultiWriter(&responseBuffer, response.Writer)
			response.Writer = httpx.NewCacheResponseWriter(response.Writer, mw)

			if err = next(ctx); err != nil {
				return err
			}

			responseBody = responseBuffer.Bytes()

			// set response to cache
			if err = distributor.Set(Context(ctx), ctx.Request(), responseBody, ttl); err != nil {
				log.
					Error(Context(ctx)).
					Err(err).
					Msg("Set cache by HTTP distributor")
			}

			return nil
		}
	}
}

func isRaw(ctx echo.Context) bool {
	return convert.Bool(Context(ctx).Value(rawResponseKey))
}
