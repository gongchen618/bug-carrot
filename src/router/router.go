package router

import (
	"bug-carrot/src/controller"
	"github.com/labstack/echo/v4"
)

func InitRouter(g *echo.Group) {
	g.POST("/reverse", controller.QQReverseHTTPMiddleHandler)
	g.POST("/hello", controller.HelloWorldHandler)
}
