package router

import (
	"bug-carrot/controller"
	"github.com/labstack/echo/v4"
)

func InitRouter(g *echo.Group) {
	g.POST("/reverse", controller.QQReverseHTTPMiddleHandler)
}
