package plugin

import (
	"bug-carrot/controller"
	"bug-carrot/controller/param"
	"bug-carrot/util"
)

type homework struct {
	PluginName string
}

func (p *homework) GetPluginName() string {
	return p.PluginName
}

func (p *homework) IsTime() bool {
	return false
}

func (p *homework) DoTime() error {
	return nil
}

func (p *homework) IsMatched(msg param.GroupMessage) bool {
	return util.IsWordInMessage("n", []string{"作业"}, msg)
}

func (p *homework) DoMatched(msg param.GroupMessage) error {
	if util.IsWordInMessage("n", []string{"微积分"}, msg) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject("微积分"))
		return nil
	}
	if util.IsWordInMessage("n", []string{"大物", "大雾", "大学物理"}, msg) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject("大物"))
		return nil
	}
	if util.IsWordInMessage("n", []string{"离散", "离散数学"}, msg) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject("离散"))
		return nil
	}
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetHomeworkStringSubject(""))
	return nil
}

func (p *homework) Listen(msg param.GroupMessage) {

}

func (p *homework) Close() {
}

func HomeworkPluginRegister() {
	p := &homework{
		PluginName: "homework",
	}
	controller.PluginRegister(p)
}
