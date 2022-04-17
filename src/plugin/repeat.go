package plugin

import (
	"bug-carrot/controller"
	"bug-carrot/controller/param"
	"bug-carrot/util"
)

type repeat struct {
	PluginName string
	RepeatStr  string
	RepeatCnt  int
}

func (p *repeat) GetPluginName() string {
	return p.PluginName
}

func (p *repeat) IsTime() bool {
	return false
}

func (p *repeat) DoTime() error {
	return nil
}

func (p *repeat) IsMatched(msg param.GroupMessage) bool {
	return false
}

func (p *repeat) DoMatched(msg param.GroupMessage) error {
	return nil
}

func (p *repeat) Listen(msg param.GroupMessage) {
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
		PluginName: "repeat",
		RepeatStr:  "",
		RepeatCnt:  0,
	}
	controller.PluginRegister(p)
}
