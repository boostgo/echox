package echox

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/boostgo/contextx"
	"github.com/boostgo/defaultx"
	"github.com/boostgo/httpx"
	"github.com/boostgo/validatex"
	"github.com/labstack/echo/v4"
)

// Param returns [param.Param] object got from named path variable or not found param error.
func Param(ctx echo.Context, paramName string) (httpx.Param, error) {
	if err := contextx.Validate(Context(ctx)); err != nil {
		return httpx.EmptyParam(), err
	}

	value := ctx.Param(paramName)
	if value == "" {
		return httpx.EmptyParam(), httpx.NewEmptyPathParamError(paramName)
	}

	return httpx.NewParam(value), nil
}

// QueryParam returns query param variable as [param.Param] object or empty [param.Param] object
// if query param is not found.
func QueryParam(ctx echo.Context, queryParamName string) httpx.Param {
	value := ctx.QueryParam(queryParamName)
	if value == "" {
		return httpx.EmptyParam()
	}

	return httpx.NewParam(value)
}

// Parse try to parse request body to provided export object (must be pointer to structure object).
//
// After success parsing request body, run format converting (for "format" tags)
//
// After success format converting, run structure validation (for "validate" tags)
func Parse(ctx echo.Context, export any) error {
	if err := contextx.Validate(Context(ctx)); err != nil {
		return err
	}

	if err := ctx.Bind(export); err != nil {
		return newParseRequestBodyError(ctx, err)
	}

	if err := validatex.Get().Struct(export); err != nil {
		return err
	}

	if err := defaultx.Set(export); err != nil {
		return err
	}

	return nil
}

// Body returns request body as []byte (slice of bytes)
func Body(ctx echo.Context) (body []byte, err error) {
	return httpx.RequestBody(ctx.Request())
}

// Context returns request context as context.Context object
func Context(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}

// SetContext sets new context to echo.Context
func SetContext(ctx echo.Context, native context.Context) {
	ctx.SetRequest(ctx.Request().WithContext(native))
}

// Set sets new key-value pair as context to request context.
func Set(ctx echo.Context, key string, value any) {
	native := Context(ctx)
	native = context.WithValue(native, key, value)
	SetContext(ctx, native)
}

// File returns file as []byte (slice of bytes) from request by file name.
//
// Request body must be form data
func File(ctx echo.Context, name string) ([]byte, error) {
	if err := contextx.Validate(Context(ctx)); err != nil {
		return nil, err
	}

	header, err := ctx.FormFile(name)
	if err != nil {
		return nil, httpx.ErrReadFormFile.SetError(err)
	}

	file, err := header.Open()
	if err != nil {
		return nil, httpx.ErrOpenFormFile.SetError(err)
	}
	defer file.Close()

	return io.ReadAll(file)
}

// ParseForm get all form data object and convert them to map with [param.Param] objects.
//
// Notice: in this map no any files. Parse them by [File] function
func ParseForm(ctx echo.Context) (map[string]httpx.Param, error) {
	if err := contextx.Validate(Context(ctx)); err != nil {
		return nil, err
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	exportMap := make(map[string]httpx.Param)
	for key, values := range form.Value {
		if len(values) == 0 {
			continue
		}

		exportMap[key] = httpx.NewParam(values[0])
	}

	return exportMap, nil
}

// Header returns request header by provided name.
func Header(ctx echo.Context, key string) httpx.Param {
	return httpx.RequestHeader(ctx.Request(), key)
}

// HeadersRaw return all headers as map with slice of values
func HeadersRaw(ctx echo.Context) map[string][]string {
	return ctx.Request().Header
}

// Headers return all headers as map with joined values
func Headers(ctx echo.Context) map[string]any {
	return httpx.RequestHeaders(ctx.Request())
}

// SetHeader sets new header to response
func SetHeader(ctx echo.Context, key, value string) {
	ctx.Response().Header().Set(key, value)
}

// Cookie returns request cookie by provided key
func Cookie(ctx echo.Context, key string) httpx.Param {
	return httpx.RequestCookie(ctx.Request(), key)
}

// CookiesRaw return all cookies as http.Cookie slice
func CookiesRaw(ctx echo.Context) []*http.Cookie {
	return ctx.Request().Cookies()
}

// Cookies return all cookies as map
func Cookies(ctx echo.Context) map[string]any {
	return httpx.RequestCookies(ctx.Request())
}

// SetCookie sets new cookie to response
func SetCookie(ctx echo.Context, key, value string, ttl ...time.Duration) {
	ctx.SetCookie(httpx.NewCookie(key, value, ttl...))
}

func StatusCode(ctx echo.Context) int {
	return ctx.Response().Status
}
