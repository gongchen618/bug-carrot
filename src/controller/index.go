package controller

import (
	"bug-carrot/util/context"
	"github.com/labstack/echo/v4"
)

// direct contact some API
func init() {

}

func HelloWorldHandler(c echo.Context) error {
	return context.Success(c, "hello world")
}
