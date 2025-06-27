package echox

import (
	"errors"
	"net/http"

	"github.com/boostgo/httpx"
	"github.com/labstack/echo/v4"
)

type httpErrorContext struct {
	Message     string `json:"message"`
	Accept      string `json:"accept"`
	ContentType string `json:"content_type"`
}

func newParseRequestBodyError(ctx echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpx.ErrParseRequestBody.
			SetError(err).
			SetData(httpErrorContext{
				Message:     err.Error(),
				Accept:      Header(ctx, "Accept").String(),
				ContentType: Header(ctx, "Content-Type").String(),
			})
	}

	return httpx.ErrParseRequestBody.SetError(err)
}

type routeNotFoundContext struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

func newRouteNotFoundError(request *http.Request) error {
	return httpx.ErrRouteNotFound.SetData(routeNotFoundContext{
		Method: request.Method,
		URL:    request.RequestURI,
	})
}
