package plugin

import (
	"bug-carrot/config"
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

// 以上四个函数在注册时被唯一调用，并以此为依据加入相应的 queue

// IsTime : 是你需要的时间吗？
func (p *_default) IsTime() bool {
	return false
}

// DoTime : 当到了你需要的时间，要做什么呢？
func (p *_default) DoTime() error {
	return nil
}

// IsMatchedGroup : 是你想收到的群 @ 消息吗？
func (p *_default) IsMatchedGroup(msg param.GroupMessage) bool {
	return true
}

// DoMatchedGroup : 收到了想收到的群 @ 消息，要做什么呢？
func (p *_default) DoMatchedGroup(msg param.GroupMessage) error {
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotGroupPuzzled)
	return nil
}

// IsMatchedPrivate : 是你想收到的私聊消息吗？
func (p *_default) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return true
}

// DoMatchedPrivate : 收到了想收到的私聊消息，要做什么呢？
func (p *_default) DoMatchedPrivate(msg param.PrivateMessage) error {
	if msg.SubType == "friend" {
		if config.C.RiskControl {
			util.QQSend(msg.UserId, constant.CarrotRiskControlAngry)
		} else {
			util.QQSend(msg.UserId, constant.CarrotFriendNotAdmin)
		}
	}
	return nil
}

// Listen : 监听到非 @ 的群消息，要做什么呢？
func (p *_default) Listen(msg param.GroupMessage) {

}

// Close : 项目要关闭了，要做什么呢？
func (p *_default) Close() {
}

// DefaultPluginRegister : 创造一个插件实例，并调用 controller.PluginRegister
// 在 main.go 的 pluginRegister 函数中调用来实现注册
func DefaultPluginRegister() {
	p := &_default{
		Index: param.PluginIndex{
			PluginName:            "default", // 插件名称
			FlagCanTime:           false,     // 是否能在特殊时间做出行为
			FlagCanMatchedGroup:   true,      // 是否能回应群聊@消息
			FlagCanMatchedPrivate: true,      // 是否能回应私聊消息
			FlagCanListen:         false,     // 是否能监听群消息
		},
	}
	controller.PluginRegister(p)
}
