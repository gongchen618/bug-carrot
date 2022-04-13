package main

import (
	"bug-carrot/config"
	"bug-carrot/router"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	e := echo.New()
	router.InitRouter(e.Group(config.C.App.Addr))
	log.Fatal(e.Start(config.C.App.Addr))
}
