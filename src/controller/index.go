package controller

import (
	"bug-carrot/src/util/context"
	"github.com/labstack/echo/v4"
)

func init() {

}

func HelloWorldHandler(c echo.Context) error {
	return context.Success(c, "hello world")
}
