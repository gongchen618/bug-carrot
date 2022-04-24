package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/model"
	param2 "bug-carrot/param"
	"bug-carrot/util"
	"fmt"
	"strings"
)

type homework struct {
	Index param2.PluginIndex
}

func (p *homework) GetPluginName() string {
	return p.Index.PluginName
}
func (p *homework) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *homework) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *homework) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *homework) CanListen() bool {
	return p.Index.FlagCanListen
}

func (p *homework) IsTime() bool {
	return false
}
func (p *homework) DoTime() error {
	return nil
}

func (p *homework) IsMatchedGroup(msg param2.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"作业"})
}
func (p *homework) DoMatchedGroup(msg param2.GroupMessage) error {
	if msg.WordsMap.ExistWord("n", []string{"微积分"}) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject("微积分"))
		return nil
	}
	if msg.WordsMap.ExistWord("n", []string{"大物", "大雾", "大学物理"}) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject("大物"))
		return nil
	}
	if msg.WordsMap.ExistWord("n", []string{"离散", "离散数学"}) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject("离散"))
		return nil
	}
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject(""))
	return nil
}

func (p *homework) IsMatchedPrivate(msg param2.PrivateMessage) bool {
	return msg.UserId == config.C.Plugin.Homework.Admin && msg.WordsMap.ExistWord("n", []string{"作业"})
}
func (p *homework) DoMatchedPrivate(msg param2.PrivateMessage) error {
	m := model.GetModel()
	defer m.Close()

	str := strings.Split(msg.RawMessage, " ")
	if len(str) >= 2 { // 格式应该是 作业 xx xx xx
		switch str[1] {
		case "delete":
			if len(str) >= 4 {
				homeworkDelete(msg.UserId, str[2], str[3])
				return nil
			}
		case "add":
			if len(str) >= 4 {
				homeworkAdd(msg.UserId, str[2], str[3])
				return nil
			}
		case "show":
			homeworkShow(msg.UserId)
			return nil
		case "clear":
			homeworkClear(msg.UserId)
			return nil
		}
	}

	util.QQSend(msg.UserId, constant.CarrotGroupPuzzled)
	return nil
}

func (p *homework) Listen(msg param2.GroupMessage) {

}

func (p *homework) Close() {
}

func HomeworkPluginRegister() {
	p := &homework{
		Index: param2.PluginIndex{
			PluginName:            "homework",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         false,
		},
	}
	controller.PluginRegister(p)
}

func homeworkDelete(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	homework := param2.Homework{
		Subject: subject,
		Context: context,
	}

	if err := m.DeleteHomework(homework); err != nil {
		util.QQSend(id, constant.CarrotHomeworkDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkDeleteSuccess)
	util.QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("删除作业 %s %s", subject, context))
}

func homeworkShow(id int64) {
	util.QQSend(id, util.GetHomeworkString())
}

func homeworkAdd(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	homework := param2.Homework{
		Subject: subject,
		Context: context,
	}

	if err := m.AddHomework(homework); err != nil {
		util.QQSend(id, constant.CarrotHomeworkAddFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkAddSuccess)
	util.QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("新增作业 %s %s", subject, context))
}

func homeworkClear(id int64) {
	m := model.GetModel()
	defer m.Close()

	if err := m.ClearAllHomework(); err != nil {
		util.QQSend(id, constant.CarrotHomeworkDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkDeleteSuccess)
	util.QQGroupSend(config.C.Plugin.Homework.Group, "新的作业正在初始化...")
}
