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

// main 是项目主入口
func main() {
	e := echo.New()
	router.InitRouter(e.Group(config.C.App.Prefix))

	pluginRegister()
	go controller.WorkTimePlugins()

	go signalWaiter()
	log.Fatal(e.Start(config.C.App.Addr))
}

// pluginRegister: 插件的总注册函数
// 新插件需要在这里调用 Register 函数
// 注意顺序
func pluginRegister() {
	plugin.HomeworkPluginRegister() // 作业
	plugin.FoodPluginRegister()     // 吃什么
	plugin.WeatherPluginRegister()  // 天气
	plugin.DicePluginRegister()     // 骰子

	plugin.GoodMorningPluginRegister() // 早安
	plugin.GoodNightPluginRegister()   // 晚安
	plugin.RepeatPluginRegister()      // 重复

	plugin.DefaultPluginRegister() // 默认回复
}

// signalWaiter 接收项目停止时的 SIGINT 和 SIGTERM 信号
// 然后调用所有插件的 close()
func signalWaiter() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	fmt.Println("exiting") //
	controller.ClosePlugins()
	os.Exit(0)
}
