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
	"time"
)

type homework struct {
	Index param2.PluginIndex
}

func (p *homework) GetPluginName() string {
	return p.Index.PluginName
}
func (p *homework) GetPluginAuthor() string {
	return p.Index.PluginAuthor
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
func (p *homework) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *homework) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
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
func (p *homework) DoMatchedGroup(msg param2.GroupMessage) error { // 还没写具体科目查询(心虚)
	if msg.WordsMap.ExistWord("n", []string{"微积分"}) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getHomeworkSubjectString("微积分"))
		return nil
	}
	if msg.WordsMap.ExistWord("n", []string{"大物", "大雾", "大学物理"}) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getHomeworkSubjectString("大物"))
		return nil
	}
	if msg.WordsMap.ExistWord("n", []string{"离散", "离散数学"}) {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getHomeworkSubjectString("离散"))
		return nil
	}
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getHomeworkString())
	return nil
}

func (p *homework) IsMatchedPrivate(msg param2.PrivateMessage) bool {
	if config.C.RiskControl && msg.WordsMap.ExistWord("n", []string{"作业"}) {
		return true
	}
	return msg.UserId == config.C.Plugin.Homework.Admin && strings.HasPrefix(msg.RawMessage, "作业")
}
func (p *homework) DoMatchedPrivate(msg param2.PrivateMessage) error { // 格式：作业 xx xx xx
	if msg.UserId == config.C.Plugin.Homework.Admin && strings.HasPrefix(msg.RawMessage, "作业") {
		str := strings.Split(msg.RawMessage, " ") // 没有考虑错误情况 因为是 admin private message
		if len(str) >= 2 {
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
	} else if config.C.RiskControl && msg.WordsMap.ExistWord("n", []string{"作业"}) {
		util.QQSend(msg.UserId, getHomeworkString())
		return nil
	}
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
			PluginAuthor:          "gongchen618",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         false,
			FlagUseDatabase:       true,
			FlagIgnoreRiskControl: false,
		},
	}
	controller.PluginRegister(p)
}

func homeworkDelete(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	hw := param2.Homework{
		Subject: subject,
		Context: context,
	}

	if err := m.DeleteHomework(hw); err != nil {
		util.QQSend(id, constant.CarrotHomeworkDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkDeleteSuccess)
	util.QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("删除作业 %s %s", subject, context))
}

func homeworkShow(id int64) {
	util.QQSend(id, getHomeworkString())
}

func homeworkAdd(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	hw := param2.Homework{
		Subject: subject,
		Context: context,
	}

	if err := m.AddHomework(hw); err != nil {
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

func getHomeworkString() string { // 暂时也还没有处理时间
	m := model.GetModel()
	defer m.Close()

	timeR := time.Now()
	d, err := time.ParseDuration("-300h")
	if err != nil {
		util.ErrorPrint(err, timeR, "time init")
		return constant.CarrotHomeworkShowFailed
	}
	timeL := timeR.Add(d)
	homeworks, err := m.GetHomeworkByTimeRange(timeL, timeR)
	if err != nil {
		util.ErrorPrint(err, timeR, "mongo")
		return constant.CarrotHomeworkShowFailed
	}

	if len(homeworks) == 0 {
		return constant.CarrotHomeworkShowEmpty
	}

	message := constant.CarrotHomeworkShowStart
	subjectMap := make(map[string]int)
	subjectInfoMap := make(map[string]string)
	for _, hw := range homeworks {
		cnt, exist := subjectMap[hw.Subject]
		info, exist := subjectInfoMap[hw.Subject]
		if exist == false {
			subjectMap[hw.Subject] = 1
			subjectInfoMap[hw.Subject] = fmt.Sprintf("(%d)%s", cnt+1, hw.Context)
		} else {
			subjectMap[hw.Subject] = cnt + 1
			subjectInfoMap[hw.Subject] = fmt.Sprintf("%s (%d)%s", info, cnt+1, hw.Context)
		}
	}
	for subject, info := range subjectInfoMap {
		message = fmt.Sprintf("%s\n【%s】%s", message, subject, info)
	}

	return message
}

func getHomeworkSubjectString(subject string) string {
	return getHomeworkString()
}
