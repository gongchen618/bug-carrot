package plugin

import (
	"bug-carrot/config"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
)

type repeat struct {
	Index     param.PluginIndex
	RepeatStr string
	RepeatCnt int
}

func (p *repeat) GetPluginName() string {
	return p.Index.PluginName
}
func (p *repeat) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *repeat) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *repeat) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *repeat) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *repeat) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *repeat) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *repeat) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

func (p *repeat) IsTime() bool {
	return false
}
func (p *repeat) DoTime() error {
	return nil
}

func (p *repeat) IsMatchedGroup(msg param.GroupMessage) bool {
	return false
}
func (p *repeat) DoMatchedGroup(msg param.GroupMessage) error {
	return nil
}

func (p *repeat) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return false
}
func (p *repeat) DoMatchedPrivate(msg param.PrivateMessage) error {
	return nil
}

func (p *repeat) Listen(msg param.GroupMessage) {
	if config.C.RiskControl {
		return
	}
	if msg.RawMessage == p.RepeatStr {
		p.RepeatCnt++
	} else {
		p.RepeatStr = msg.RawMessage
		p.RepeatCnt = 0
	}
	if p.RepeatCnt == 4 {
		util.QQGroupSend(msg.GroupId, p.RepeatStr)
	}
}

func (p *repeat) Close() {
}

func RepeatPluginRegister() {
	p := &repeat{
		Index: param.PluginIndex{
			PluginName:            "repeat",
			PluginAuthor:          "gongchen618",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   false,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         !config.C.RiskControl,
			FlagUseDatabase:       false,
			FlagIgnoreRiskControl: false,
		},
		RepeatStr: "",
		RepeatCnt: 0,
	}
	controller.PluginRegister(p)
}
