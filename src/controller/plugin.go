package controller

import (
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
	if p.CanTime() {
		PluginTime = append(PluginTime, p)
	}
	if p.CanMatchedGroup() {
		PluginGroup = append(PluginGroup, p)
	}
	if p.CanMatchedPrivate() {
		PluginPrivate = append(PluginPrivate, p)
	}
	if p.CanListen() {
		PluginListen = append(PluginListen, p)
	}
	Plugin = append(Plugin, p)
	fmt.Println(fmt.Sprintf("插件 %s 已装载", p.GetPluginName()))
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
