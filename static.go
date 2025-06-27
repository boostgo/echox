package echox

import (
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Swagger(path string) {
	GET(path, echoSwagger.WrapHandler)
}
