package plugin

import (
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
)

type _default struct {
	Index param.PluginIndex
}

func (p *_default) GetPluginName() string {
	return p.Index.PluginName
}
func (p *_default) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *_default) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *_default) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *_default) CanListen() bool {
	return p.Index.FlagCanListen
}

func (p *_default) IsTime() bool {
	return false
}
func (p *_default) DoTime() error {
	return nil
}

func (p *_default) IsMatchedGroup(msg param.GroupMessage) bool {
	return true
}
func (p *_default) DoMatchedGroup(msg param.GroupMessage) error {
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotGroupPuzzled)
	return nil
}

func (p *_default) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return true
}
func (p *_default) DoMatchedPrivate(msg param.PrivateMessage) error {
	if msg.SubType == "friend" {
		util.QQSend(msg.UserId, constant.CarrotFriendNotAdmin)
	}
	return nil
}

func (p *_default) Listen(msg param.GroupMessage) {

}

func (p *_default) Close() {
}

func DefaultPluginRegister() {
	p := &_default{
		Index: param.PluginIndex{
			PluginName:            "default",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         false,
		},
	}
	controller.PluginRegister(p)
}
