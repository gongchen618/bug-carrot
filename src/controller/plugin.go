package controller

import (
	"bug-carrot/controller/param"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	Plugin []param.PluginInterface
)

func PluginRegister(p param.PluginInterface) {
	Plugin = append(Plugin, p)
}

func MessagePluginCenter(msg param.GroupMessage) {
	for _, p := range Plugin {
		if p.IsMatched(msg) {
			if err := p.DoMatched(msg); err != nil {
				logrus.WithFields(logrus.Fields{"err": err, "plugin": p.GetPluginName(), "type": "message"}).Warn("plugin error")
			}
			return
		}
	}
}

func TimePluginCenter() {
	for {
		for _, p := range Plugin {
			if p.IsTime() {
				if err := p.DoTime(); err != nil {
					logrus.WithFields(logrus.Fields{"err": err, "plugin": p.GetPluginName(), "type": "time"}).Warn("plugin error")
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func ListenPluginCenter(msg param.GroupMessage) {
	for _, p := range Plugin {
		p.Listen(msg)
	}
}

func PluginClose() {
	for _, p := range Plugin {
		p.Close()
	}
}
