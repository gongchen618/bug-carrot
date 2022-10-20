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
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：身份认证失败，你是...？"))
		return nil
	}

	help := msg.ExistWord("v", []string{"帮助"})
	if help {
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("\"填写 帮助\"：查看本消息\n"+
			"\"填写 查看\"：查看当前未填写的收集\n"+
			"\"填写 查看 all\"：查看当前所有收集\n"+
			"\"填写 [title] [answer]\"：填写名为 [title] 的收集，答复为 [answer]"))
		return nil
	}

	ask := msg.ExistWord("v", []string{"查看"})
	if ask {
		allBallot, err := m.GetAllBallot()
		if err != nil {
			util.QQSend(msg.UserId, "FAILED：有某种不可抗力量抑制了卡洛的魔力！")
			return nil
		}
		message := "查询结果如下：\n"
		flagAll := msg.ExistWord("eng", []string{"all"})
		for _, bt := range allBallot {
			for _, mb := range bt.TargetMember {
				if mb.People.QQ == msg.UserId && (mb.AnsweredFlag == false || flagAll) {
					if mb.AnsweredFlag == false {
						message = fmt.Sprintf("%s【%s】（未回复）\n备注：%s\n", message, bt.Title, bt.Remark)
					} else {
						message = fmt.Sprintf("%s【%s】回复：%s\n备注：%s\n", message, bt.Title, mb.Answer, bt.Remark)
					}
				}
			}
		}
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage(message))
		return nil
	}

	word := strings.Split(msg.RawMessage, " ")

	if len(word) >= 3 {
		bt, err := m.GetOneBallotByTitle(word[1])
		if err != nil {
			util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage(fmt.Sprintf("FAILED：没有找到收集【%s】", word[1])))
			return nil
		}
		visFlag := false
		for _, person := range bt.TargetMember {
			if person.People.QQ == msg.UserId {
				visFlag = true
			}
		}
		if visFlag == false {
			util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage(fmt.Sprintf("FAILED：收集【%s】好像不需要你回答哦", word[1])))
			return nil
		}
		_, err = m.UpdateAnswerForOneMember(word[1], strings.Join(word[2:], " "), member.Name)
		if err != nil {
			util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：没回复成功，我的。"))
			return nil
		}
		util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("ACCEPT: 回答已接收！"))
		return nil
	}

	util.QQSend(msg.UserId, util.GetHitokotoWarpedMessage("FAILED：解析失败，你肯定又在说怪话了。\n没有的话，要不要问问我 \"填写 帮助\" 呢？"))
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
