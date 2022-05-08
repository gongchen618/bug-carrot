package controller

import (
	"bug-carrot/config"
	"bug-carrot/param"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	PluginTime    []param.PluginInterface
	PluginGroup   []param.PluginInterface
	PluginPrivate []param.PluginInterface
	PluginListen  []param.PluginInterface
	Plugin        []param.PluginInterface
)

func PluginRegister(p param.PluginInterface) {
	if p.NeedDatabase() && !config.C.DatabaseUse {
		fmt.Println(fmt.Sprintf("数据库缺失，插件 [%s]%s 未装载", p.GetPluginAuthor(), p.GetPluginName()))
		return
	}
	if p.CanTime() {
		PluginTime = append(PluginTime, p)
	}
	if p.CanMatchedGroup() {
		if !config.C.RiskControl || p.DoIgnoreRiskControl() {
			PluginGroup = append(PluginGroup, p)
		}
	}
	if p.CanMatchedPrivate() {
		PluginPrivate = append(PluginPrivate, p)
	}
	if p.CanListen() {
		if !config.C.RiskControl || p.DoIgnoreRiskControl() {
			PluginListen = append(PluginListen, p)
		}
	}

	Plugin = append(Plugin, p)

	if config.C.RiskControl && !p.DoIgnoreRiskControl() {
		fmt.Println(fmt.Sprintf("当前位于风控场景，插件 [%s]%s 的群聊和监听功能未装载", p.GetPluginAuthor(), p.GetPluginName()))
	} else {
		fmt.Println(fmt.Sprintf("插件 [%s]%s 已装载", p.GetPluginAuthor(), p.GetPluginName()))
	}
}

func WorkGroupMessagePlugins(msg param.GroupMessage) {
	for _, p := range PluginGroup {
		if p.IsMatchedGroup(msg) {
			if err := p.DoMatchedGroup(msg); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "plugin": p.GetPluginName(), "type": "message"}).Warn("plugin error")
			}
			return
		}
	}
}

func WorkPrivateMessagePlugins(msg param.PrivateMessage) {
	for _, p := range PluginPrivate {
		if p.IsMatchedPrivate(msg) {
			if err := p.DoMatchedPrivate(msg); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "plugin": p.GetPluginName(), "type": "message"}).Warn("plugin error")
			}
			return
		}
	}
}

func WorkTimePlugins() { // 10s 一次轮询
	for {
		for _, p := range PluginTime {
			if p.IsTime() {
				if err := p.DoTime(); err != nil {
					logrus.WithFields(logrus.Fields{"err": err, "plugin": p.GetPluginName(), "type": "time"}).Warn("plugin error")
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func WorkListenPlugins(msg param.GroupMessage) {
	for _, p := range PluginListen {
		p.Listen(msg)
	}
}

func ClosePlugins() {
	for _, p := range Plugin {
		p.Close()
	}
}
