package plugin

import (
	"bug-carrot/controller"
	"bug-carrot/model"
	"bug-carrot/param"
	"bug-carrot/util"
	"fmt"
	"strings"
)

type ballot struct {
	Index param.PluginIndex
}

func (p *ballot) GetPluginName() string {
	return p.Index.PluginName
}
func (p *ballot) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *ballot) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *ballot) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *ballot) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *ballot) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *ballot) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *ballot) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

// 以上六个函数在注册时被唯一调用，并以此为依据加入相应的 queue
// 无需修改

// IsTime : 是你需要的时间吗？
func (p *ballot) IsTime() bool {
	return false
}

// DoTime : 当到了你需要的时间，要做什么呢？
func (p *ballot) DoTime() error {
	return nil
}

// IsMatchedGroup : 是你想收到的群 @ 消息吗？
func (p *ballot) IsMatchedGroup(msg param.GroupMessage) bool {
	return false
}

// DoMatchedGroup : 收到了想收到的群 @ 消息，要做什么呢？
func (p *ballot) DoMatchedGroup(msg param.GroupMessage) error {
	return nil
}

// IsMatchedPrivate : 是你想收到的私聊消息吗？
func (p *ballot) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return len(msg.RawMessage) >= len("填写") && msg.RawMessage[:len("填写")] == "填写"
}

// DoMatchedPrivate : 收到了想收到的私聊消息，要做什么呢？
// 备注：我们建议大部分功能只对群聊开启，增强 bot 在群聊中的存在感，私聊功能可以提供给管理员
func (p *ballot) DoMatchedPrivate(msg param.PrivateMessage) error {
	if msg.SubType != "friend" {
		return nil
	}

	m := model.GetModel()
	defer m.Close()

	member, err := m.GetOneFamilyMemberByQQ(msg.UserId)
	if err != nil {
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：身份认证失败，你是谁？"))
		return nil
	}

	help := msg.ExistWord("v", []string{"帮助"})
	if help {
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("\"填写 帮助\"：查看本消息\n"+
			"\"填写 查询\"：查看当前收集\n"+
			"\"填写 [title] [option]\"：填写名为 [title] 的收集，并选择 [option]"))
		return nil
	}

	ask := msg.ExistWord("v", []string{"查询"})
	if ask {
		allBallot, err := m.GetAllBallot()
		if err != nil {
			util.QQSend(msg.UserId, "FAILED：有某种不可抗力量抑制了卡洛的魔力！")
			return nil
		}
		message := ""
		for _, bt := range allBallot {
			for _, mb := range bt.TargetMember {
				if mb.Info.QQ == msg.UserId {
					message = fmt.Sprintf("%s你在【%s】中的选择为：%s选项包括：%s\n", message, bt.Title, mb.Option, bt.OfferedOptions)
				}
			}
		}
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage(message))
		return nil
	}

	word := strings.Split(msg.RawMessage, " ")

	if len(word) == 3 {
		bt, err := m.GetOneBallotByTitle(word[1])
		vis := false
		for _, option := range bt.OfferedOptions {
			if option == word[2] {
				vis = true
			}
		}
		if err != nil || !vis {
			util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：选择未成功，你的输入非法。"))
			return nil
		}
		_, err = m.UpdateOptionsOnOneBallotForMembers(word[1], word[2], []string{member.Name})
		if err != nil {
			util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：选择未成功，我的。"))
			return nil
		}
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("ACCEPT: 选项已接收！"))
		return nil
	}

	util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：解析失败，你肯定又在说怪话了。"))
	return nil
}

func (p *ballot) Listen(msg param.GroupMessage) {}

func (p *ballot) Close() {
}

func BallotPluginRegister() {
	p := &ballot{
		Index: param.PluginIndex{
			PluginName:            "ballot",      // 插件名称
			PluginAuthor:          "gongchen618", // 插件作者
			FlagCanTime:           false,         // 是否能在特殊时间做出行为
			FlagCanMatchedGroup:   false,         // 是否能回应群聊@消息
			FlagCanMatchedPrivate: true,          // 是否能回应私聊消息
			FlagCanListen:         false,         // 是否能监听群消息
			FlagUseDatabase:       true,          // 是否用到了数据库（配置文件中配置不使用数据库的话，用到了数据库的插件会不运行）
			FlagIgnoreRiskControl: false,         // 是否无视风控（为 true 且 RiskControl=true 时将自动无视群聊功能，建议设置为 false）
		},
	}
	controller.PluginRegister(p)
}
