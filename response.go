package echox

import "github.com/labstack/echo/v4"

type FailureMiddleware func(ctx echo.Context, statusCode int, err error)

var failureMiddlewares []FailureMiddleware

func RegisterFailureMiddleware(m FailureMiddleware) {
	failureMiddlewares = append(failureMiddlewares, m)
}
