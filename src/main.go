package main

import (
	"bug-carrot/config"
	"bug-carrot/controller"
	"bug-carrot/plugin"
	"bug-carrot/router"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	e := echo.New()
	router.InitRouter(e.Group(config.C.App.Addr))

	pluginRegister()
	go controller.TimePluginCenter()

	go signalWaiter()
	log.Fatal(e.Start(config.C.App.Addr))
}

func pluginRegister() { // 总注册函数，注意顺序
	plugin.HomeworkPluginRegister()
	plugin.WeatherPluginRegister()
	plugin.GoodMorningPluginRegister()
	plugin.GoodNightPluginRegister()
	plugin.DefaultPluginRegister() // 默认回复
}

func signalWaiter() { // 接收 SIGINT 和 SIGTERM 信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	fmt.Println("exiting") //
	controller.PluginClose()
	os.Exit(0)
}
