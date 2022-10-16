package router

import (
	"bug-carrot/controller"
	"github.com/labstack/echo/v4"
)

func InitRouter(g *echo.Group) {
	g.POST("/reverse", controller.QQReverseHTTPMiddleHandler)
	g.POST("/hello", controller.HelloWorldHandler)

	initFamilyAPIRouter(g.Group("/family"))
	initMusterAPIRouter(g.Group("/muster"))
	initBallotAPIRouter(g.Group("/ballot"))
}

func initFamilyAPIRouter(g *echo.Group) {
	g.POST("", controller.CreateOneFamilyMemberRequestHandler)
	g.GET("/all", controller.GetAllFamilyMembersRequestHandler)
	g.DELETE("", controller.DeleteOneFamilyMemberByStudentIDRequestHandler)
	g.PUT("", controller.UpdateOneFamilyMemberRequestHandler)
}

func initMusterAPIRouter(g *echo.Group) {
	g.GET("/all", controller.GetAllMusterRequestHandler)
	g.POST("", controller.CreateOneMusterByNameRequestHandler)
	g.DELETE("", controller.DeleteOneMusterByNameRequestHandler)
	g.POST("/people", controller.AddPersonsToOneMusterRequestHandler)
	g.DELETE("/people", controller.DeletePersonsOnOneMusterRequestHandler)
}

func initBallotAPIRouter(g *echo.Group) {
	g.GET("/all", controller.GetAllBallotRequestHandler)
	g.POST("", controller.CreateOneBallotByTitleRequestHandler)
	g.DELETE("", controller.DeleteOneBallotByTitleRequestHandler)
	g.POST("/option", controller.AddAnOptionToOneBallotRequestHandler)
	g.DELETE("/option", controller.DeleteAnOptionOnOneBallotRequestHandler)
	g.PUT("/member", controller.UpdateOptionsOnOneBallotForMembersRequestHandler)
}
