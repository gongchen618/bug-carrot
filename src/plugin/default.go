package plugin

import (
	"bug-carrot/src/config"
	"bug-carrot/src/constant"
	"bug-carrot/src/controller"
	"bug-carrot/src/param"
	"bug-carrot/src/util"
)

type _default struct {
	Index param.PluginIndex
}

func (p *_default) GetPluginName() string {
	return p.Index.PluginName
}
func (p *_default) GetPluginAuthor() string {
	return p.Index.PluginAuthor
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
func (p *_default) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *_default) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

// 以上六个函数在注册时被唯一调用，并以此为依据加入相应的 queue
// 无需修改

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
	if !config.C.RiskControl { // 备注：这里 RiskControl 的全局变量用于风控
		// 在 bot 运行已经稳定的场景下，我们认为新增的插件将此变量忽略掉是可以接受的
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotGroupPuzzled)
	}
	return nil
}

// IsMatchedPrivate : 是你想收到的私聊消息吗？
func (p *_default) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return true
}

// DoMatchedPrivate : 收到了想收到的私聊消息，要做什么呢？
// 备注：我们建议大部分功能只对群聊开启，增强 bot 在群聊中的存在感，私聊功能可以提供给管理员
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
// 备注：我们建议只对极少数的功能采取监听行为
// 除去整活效果较好的特殊场景，我们一般希望 bot 只有在被 @ 到的时候才会对应s发言
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
			PluginName:            "default",     // 插件名称
			PluginAuthor:          "gongchen618", // 插件作者
			FlagCanTime:           false,         // 是否能在特殊时间做出行为
			FlagCanMatchedGroup:   true,          // 是否能回应群聊@消息
			FlagCanMatchedPrivate: true,          // 是否能回应私聊消息
			FlagCanListen:         false,         // 是否能监听群消息
			FlagUseDatabase:       false,         // 是否用到了数据库（配置文件中配置不使用数据库的话，用到了数据库的插件会不运行）
			FlagIgnoreRiskControl: true,          // 是否无视风控（为 true 且 RiskControl=true 时将自动无视群聊功能，建议设置为 false）
		},
	}
	controller.PluginRegister(p)
}
