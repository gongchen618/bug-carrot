package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
)

type vjudgeRank struct {
	Index param.PluginIndex
}

func (p *vjudgeRank) GetPluginName() string {
	return p.Index.PluginName
}
func (p *vjudgeRank) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *vjudgeRank) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *vjudgeRank) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *vjudgeRank) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *vjudgeRank) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *vjudgeRank) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *vjudgeRank) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

// 以上六个函数在注册时被唯一调用，并以此为依据加入相应的 queue
// 无需修改

// IsTime : 是你需要的时间吗？
func (p *vjudgeRank) IsTime() bool {
	return false
}

// DoTime : 当到了你需要的时间，要做什么呢？
func (p *vjudgeRank) DoTime() error {
	return nil
}

// IsMatchedGroup : 是你想收到的群 @ 消息吗？
func (p *vjudgeRank) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"榜单"})
}

// DoMatchedGroup : 收到了想收到的群 @ 消息，要做什么呢？

var (
	contestIDDiv1 = "519675"
	contestIDDiv2 = "519676"
)

func (p *vjudgeRank) DoMatchedGroup(msg param.GroupMessage) error {
	var message string
	if msg.WordsMap.ExistWord("eng", []string{"div2"}) {
		message = util.GetRankString(contestIDDiv2)
	} else {
		message = util.GetRankString(contestIDDiv1)
	}
	util.QQGroupSend(msg.GroupId, message)
	return nil
}

func (p *vjudgeRank) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return true
}

func (p *vjudgeRank) DoMatchedPrivate(msg param.PrivateMessage) error {
	if msg.SubType == "friend" {
		if config.C.RiskControl {
			util.QQSend(msg.UserId, constant.CarrotRiskControlAngry)
		} else {
			util.QQSend(msg.UserId, constant.CarrotFriendNotAdmin)
		}
	}
	return nil
}

func (p *vjudgeRank) Listen(msg param.GroupMessage) {

}

func (p *vjudgeRank) Close() {
}

func VjudgeRankPluginRegister() {
	p := &vjudgeRank{
		Index: param.PluginIndex{
			PluginName:            "vjudge_rank", // 插件名称
			PluginAuthor:          "Smokey_Days", // 插件作者
			FlagCanTime:           true,          // 是否能在特殊时间做出行为
			FlagCanMatchedGroup:   true,          // 是否能回应群聊@消息
			FlagCanMatchedPrivate: false,         // 是否能回应私聊消息
			FlagCanListen:         false,         // 是否能监听群消息
			FlagUseDatabase:       false,         // 是否用到了数据库（配置文件中配置不使用数据库的话，用到了数据库的插件会不运行）
			FlagIgnoreRiskControl: false,         // 是否无视风控（为 true 且 RiskControl=true 时将自动无视群聊功能，建议设置为 false）
		},
	}
	controller.PluginRegister(p)
}
