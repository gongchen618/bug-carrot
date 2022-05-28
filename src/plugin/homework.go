package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/model"
	param2 "bug-carrot/param"
	"bug-carrot/util"
	"errors"
	"fmt"
	"strings"
	"time"
)

type homework struct {
	Index      param2.PluginIndex
	LasWeekday time.Weekday
}

var (
	errorWeekdayParseFailed = errors.New("weekday parse failed")
	errorMongoRunFailed     = errors.New("mongo run failed")
	errorNoHomework         = errors.New("no homework is found")
)

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
	weekday := time.Now().Weekday()
	if weekday == p.LasWeekday {
		return false
	}
	p.LasWeekday = weekday

	m := model.GetModel()
	defer m.Close()

	homeworks, err := m.GetHomeWorkByWeekDay(weekday)
	if err != nil {
		util.ErrorPrint(err, nil, "mongo")
		return false
	}
	if len(homeworks) == 0 {
		return false
	}

	return true
}
func (p *homework) DoTime() error {
	m := model.GetModel()
	defer m.Close()

	homeworks, err := m.GetHomeWorkByWeekDay(time.Now().Weekday())
	if err != nil {
		util.ErrorPrint(err, nil, "mongo")
		return err
	}
	if len(homeworks) == 0 {
		return errors.New("is time but no documents")
	}

	message := fmt.Sprintf("今天有作业要交哦~%s", parseHomeworkToString(homeworks))

	util.QQGroupSend(config.C.Plugin.Homework.Group, message)
	return nil
}

func (p *homework) IsMatchedGroup(msg param2.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"作业"})
}
func (p *homework) DoMatchedGroup(msg param2.GroupMessage) error {
	str := strings.Split(msg.RawMessage, " ")
	for i, word := range str {
		if word == "-t" {
			if i+1 == len(str) {
				util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), "星期解析失败!")
				return nil
			}
			message, err := getHomeworkStringByWeekDay(str[i+1])
			if err != nil {
				if err == errorNoHomework {
					util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), fmt.Sprintf("%s没有要交的作业哦!", str[i+1]))
					return nil
				}
				util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), fmt.Sprintf("星期解析失败:%s", err.Error()))
				return nil
			}
			util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), message)
			return nil
		}
	}

	message, err := getHomeworkStringAll()
	if err != nil {
		if err == errorNoHomework {
			util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), "现在没有要交的作业哦~大家一起摸鱼吧!")
			return nil
		}
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), fmt.Sprintf("作业获取失败:%s", err.Error()))
		return nil
	}
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), message)
	return nil
}

func (p *homework) IsMatchedPrivate(msg param2.PrivateMessage) bool {
	return msg.UserId == config.C.Plugin.Homework.Admin && strings.HasPrefix(msg.RawMessage, "作业")
}
func (p *homework) DoMatchedPrivate(msg param2.PrivateMessage) error { // 格式：作业 xx xx xx
	str := strings.Split(msg.RawMessage, " ")
	if len(str) >= 2 {
		switch str[1] {
		case "delete":
			if len(str) >= 4 {
				homeworkDelete(msg.UserId, str[2], strings.Join(str[3:], " "))
				return nil
			}
		case "add":
			if len(str) >= 5 {
				homeworkAdd(msg.UserId, str[2], str[3], strings.Join(str[4:], " "))
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
			PluginAuthor:          "gongchen618",
			FlagCanTime:           true,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         false,
			FlagUseDatabase:       true,
			FlagIgnoreRiskControl: false,
		},
		LasWeekday: time.Now().Weekday(),
	}
	controller.PluginRegister(p)
}

func homeworkDelete(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	if err := m.DeleteHomework(subject, context); err != nil {
		util.QQSend(id, constant.CarrotHomeworkDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkDeleteSuccess)
	util.QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("删除作业 %s %s", subject, context))
}

func homeworkShow(id int64) {
	message, err := getHomeworkStringAll()
	if err != nil {
		util.QQSend(id, fmt.Sprintf("%s%s", constant.CarrotHomeworkShowFailed, err.Error()))
		return
	}
	util.QQSend(id, message)
}

func homeworkAdd(id int64, subject string, weekday string, context string) {
	m := model.GetModel()
	defer m.Close()

	weekdayTime, err := parseStringToWeekday(weekday)
	if err != nil {
		util.QQSend(id, constant.CarrotHomeworkDeleteFailed)
		return
	}

	hw := param2.Homework{
		Subject: subject,
		Context: context,
		Weekday: weekdayTime,
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

func getHomeworkStringAll() (string, error) {
	m := model.GetModel()
	defer m.Close()

	homeworks, err := m.GetHomeworkFromNow()
	if err != nil {
		util.ErrorPrint(err, nil, "mongo")
		return "", errorMongoRunFailed
	}
	if len(homeworks) == 0 {
		return "", errorNoHomework

	}

	message := fmt.Sprintf("%s%s", constant.CarrotHomeworkShowStart, parseHomeworkToString(homeworks))

	return message, nil
}

func getHomeworkStringByWeekDay(weekday string) (string, error) {
	weekdayTime, err := parseStringToWeekday(weekday)
	if err != nil {
		return "", errorWeekdayParseFailed
	}

	m := model.GetModel()
	defer m.Close()

	homeworks, err := m.GetHomeWorkByWeekDay(weekdayTime)
	if err != nil {
		util.ErrorPrint(err, nil, "mongo")
		return "", errorMongoRunFailed
	}
	if len(homeworks) == 0 {
		return "", errorNoHomework
	}

	message := fmt.Sprintf("%s要交的作业包括：%s", weekday, parseHomeworkToString(homeworks))

	return message, nil
}

// parseHomeworkToString: WARNING: 自带一个行首换行
func parseHomeworkToString(homeworks []param2.Homework) string {
	message := ""
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

func parseStringToWeekday(weekday string) (time.Weekday, error) {
	switch weekday {
	case "Mon", "Monday", "周一", "礼拜一":
		return time.Monday, nil
	case "Tue", "Tuesday", "周二", "礼拜二":
		return time.Tuesday, nil
	case "Wed", "Wednesday", "周三", "礼拜三":
		return time.Wednesday, nil
	case "Thu", "Thursday", "周四", "礼拜四":
		return time.Thursday, nil
	case "Fri", "Friday", "周五", "礼拜五":
		return time.Friday, nil
	case "Sat", "Saturday", "周六", "礼拜六":
		return time.Saturday, nil
	case "Sun", "Sunday", "周天", "周日", "礼拜天":
		return time.Sunday, nil
	}
	return time.Sunday, errors.New("parse failed")
}
