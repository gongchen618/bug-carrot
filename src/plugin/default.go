package plugin

import (
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/controller/param"
	"bug-carrot/util"
)

type _default struct {
	PluginName string
}

func (p *_default) GetPluginName() string {
	return p.PluginName
}

func (p *_default) IsTime() bool {
	return false
}

func (p *_default) DoTime() error {
	return nil
}

func (p *_default) IsMatched(msg param.GroupMessage) bool {
	return true
}

func (p *_default) DoMatched(msg param.GroupMessage) error {
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotGroupPuzzled)
	return nil
}

func (p *_default) Listen(msg param.GroupMessage) {

}

func (p *_default) Close() {
}

func DefaultPluginRegister() {
	p := &_default{
		PluginName: "default",
	}
	controller.PluginRegister(p)
}
