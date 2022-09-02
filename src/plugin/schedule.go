package plugin

import (
	"bug-carrot/src/config"
	"bug-carrot/src/constant"
	"bug-carrot/src/controller"
	"bug-carrot/src/model"
	"bug-carrot/src/param"
	"bug-carrot/src/util"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strconv"
	"strings"
	"time"
)

type schedule struct {
	Index     param.PluginIndex
	lasQuater int
}

func (p *schedule) GetPluginName() string {
	return p.Index.PluginName
}
func (p *schedule) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *schedule) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *schedule) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *schedule) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *schedule) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *schedule) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *schedule) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

func (p *schedule) IsTime() bool {
	if time.Now().Weekday() == time.Monday {
		return true
	} // 每周一发一次

	if p.lasQuater == time.Now().Minute()/15 {
		return false
	} // 每 15 分钟检查一次近期约定
	p.lasQuater = time.Now().Minute() / 15

	m := model.GetModel()
	defer m.Close()

	schedules, err := m.GetScheduleAllFromNow()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		util.ErrorPrint(err, nil, "mongo")
		return false
	}

	for _, object := range schedules {
		duration := time.Now().Sub(object.Date.Local()).Hours()
		if math.Abs(duration) < 24 && object.Alarm24hFlag == false {
			return true
		} else if math.Abs(duration) < 1 && object.Alarm1hFlag == false {
			return true
		}
	}

	return false
}

func (p *schedule) DoTime() error {
	if time.Now().Weekday() == time.Monday {
		util.QQSend(config.C.Plugin.Schedule.Group, getScheduleStringAllFromNow())
		return nil
	}
	util.QQGroupSend(config.C.Plugin.Schedule.Group, getScheduleStringRecent())
	return nil
}

func (p *schedule) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"任务"}) ||
		msg.WordsMap.ExistWord("eng", []string{"TODO"}) ||
		msg.WordsMap.ExistWord("a", []string{"清单"}) ||
		msg.WordsMap.ExistWord("v", []string{"约定"})
}

func (p *schedule) DoMatchedGroup(msg param.GroupMessage) error { // #asd##asd#(page) // 分页还没有写完
	var keyword string
	for i, c := range msg.RawMessage {
		if c == '#' && i+2 < len(msg.RawMessage) {
			for j, c2 := range msg.RawMessage[i+1:] {
				if c2 == '#' {
					keyword = msg.RawMessage[i+1 : i+j+1]
					break
				}
			}
			break
		}
	}

	var pageString string
	for i, c := range msg.RawMessage {
		if c == 'p' && i+2 < len(msg.RawMessage) {
			for j, c2 := range msg.RawMessage[i+1:] {
				if c2 == 'p' {
					pageString = msg.RawMessage[i+1 : i+j+1]
					break
				}
			}
			break
		}
	}
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 0
	}

	var limitString string
	for i, c := range msg.RawMessage {
		if c == 'l' && i+2 < len(msg.RawMessage) {
			for j, c2 := range msg.RawMessage[i+1:] {
				if c2 == 'l' {
					limitString = msg.RawMessage[i+1 : i+j+1]
					break
				}
			}
			break
		}
	}
	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		limit = 10
	}

	if keyword == "" {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getScheduleStringAllFromNow())
		return nil
	}

	_, err = strconv.ParseInt(keyword, 10, 64)
	if err != nil {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getScheduleStringByTitleFromNow(keyword, page, limit))
		return nil
	}

	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getScheduleDetailStringById(keyword))
	return nil
}

func (p *schedule) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return msg.UserId == config.C.Plugin.Schedule.Admin && strings.HasPrefix(msg.RawMessage, "约定")
}

func (p *schedule) DoMatchedPrivate(msg param.PrivateMessage) error { // 格式：约定 func time title description
	str := strings.Split(msg.RawMessage, " ") // 没有考虑错误情况 因为是 admin private message
	if len(str) >= 2 {
		switch str[1] {
		case "delete": // delete id1 id2 id3 id4 id5 批量删除
			if len(str) >= 3 {
				for _, id := range (str)[2:] {
					scheduleDeletePrivate(msg.UserId, id)
				}
				return nil
			}
		case "add": // add time title description
			if len(str) >= 5 {
				scheduleAddPrivate(msg.UserId, str[2], str[3], strings.Join(str[4:], " "))
				return nil
			}
		case "update": // update id time title description
			if len(str) >= 6 {
				scheduleUpdatePrivate(msg.UserId, str[2], str[3], str[4], strings.Join(str[5:], " "))
				return nil
			}
		case "show": // show (title) page limit
			if len(str) == 5 {
				page, err := strconv.ParseInt(str[3], 10, 64)
				if err != nil {
					page = 0
				}
				limit, err := strconv.ParseInt(str[4], 10, 64)
				if err != nil {
					limit = 10
				}
				util.QQSend(msg.UserId, getScheduleStringByTitleFromNow(str[2], page, limit))
				return nil
			}
			util.QQSend(msg.UserId, getScheduleStringAllFromNow())
			return nil
		case "detail": // detail id
			if len(str) == 3 {
				util.QQSend(msg.UserId, getScheduleDetailStringById(str[2]))
				return nil
			}
		}
	}

	util.QQSend(msg.UserId, constant.CarrotGroupPuzzled)
	return nil
}

func (p *schedule) Listen(msg param.GroupMessage) {

}

func (p *schedule) Close() {
}

var (
	timePattern        = "2006/01/02/15:04" // time parse pattern
	timePatternSimple  = "1-2 Mon"
	timePatternComplex = "01月02日15:04,Monday"
)

func SchedulePluginRegister() {
	p := &schedule{
		Index: param.PluginIndex{
			PluginName:            "schedule",    // 插件名称
			PluginAuthor:          "gongchen618", // 插件作者
			FlagCanTime:           true,          // 是否能在特殊时间做出行为
			FlagCanMatchedGroup:   true,          // 是否能回应群聊@消息
			FlagCanMatchedPrivate: true,          // 是否能回应私聊消息
			FlagCanListen:         false,         // 是否能监听群消息
			FlagUseDatabase:       true,          // 是否用到了数据库（配置文件中配置不使用数据库的话，用到了数据库的插件会不运行）
			FlagIgnoreRiskControl: false,         // 是否无视风控（为 true 且 RiskControl=true 时将自动无视群聊功能，建议设置为 false）
		},
		lasQuater: -1,
	}
	controller.PluginRegister(p)
}

func scheduleDeletePrivate(id int64, ScheduleId string) {
	m := model.GetModel()
	defer m.Close()

	object, err := m.DeleteScheduleById(ScheduleId)
	if err != nil {
		util.QQSend(id, constant.CarrotScheduleDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotScheduleDeleteSuccess)
	util.QQGroupSend(config.C.Plugin.Schedule.Group, fmt.Sprintf("约定「%s」被废弃了!", object.Title))
}

func scheduleAddPrivate(id int64, dateStr string, title string, description string) {
	m := model.GetModel()
	defer m.Close()

	date, err := time.ParseInLocation(timePattern, dateStr, time.Local)
	if err != nil {
		util.QQSend(id, constant.CarrotScheduleAddFailed)
		util.ErrorPrint(err, nil, "date")
		return
	}

	cnt, err := m.GetScheduleCount()
	if err != nil {
		util.QQSend(id, constant.CarrotScheduleAddFailed)
		util.ErrorPrint(err, nil, "count")
		return
	}
	object := param.Schedule{
		ScheduleId:   strconv.FormatInt(cnt, 10),
		Date:         date,
		Title:        title,
		Description:  description,
		ExistFlag:    true,
		Alarm1hFlag:  false,
		Alarm24hFlag: false,
	}

	if err = m.AddSchedule(object); err != nil {
		util.QQSend(id, constant.CarrotScheduleAddFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotScheduleAddSuccess)
	util.QQGroupSend(config.C.Plugin.Schedule.Group, fmt.Sprintf("约定好了，要和卡洛在「%s」一起「%s」哦！", object.Date.Local().Format(timePatternSimple), title))
}

func scheduleUpdatePrivate(id int64, ScheduleId string, dateStr string, title string, description string) {
	m := model.GetModel()
	defer m.Close()

	date, err := time.ParseInLocation(timePattern, dateStr, time.Local)
	if err != nil {
		util.QQSend(id, constant.CarrotScheduleAddFailed)
		util.ErrorPrint(err, nil, "date")
		return
	}
	object := param.Schedule{
		Date:         date,
		Title:        title,
		Description:  description,
		ExistFlag:    true,
		Alarm1hFlag:  false,
		Alarm24hFlag: false,
	}

	updatedObject, err := m.UpdateScheduleById(ScheduleId, object)
	if err != nil {
		util.QQSend(id, constant.CarrotScheduleAddFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotScheduleUpdatedSuccess)
	util.QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("约定「%s」被更新了~", updatedObject.Title))
}

func getScheduleStringAllFromNow() string {
	m := model.GetModel()
	defer m.Close()

	schedules, err := m.GetScheduleAllFromNow()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return constant.CarrotScheduleNotFound
		}
		util.ErrorPrint(err, nil, "mongo")
		return constant.CarrotScheduleShowFailed
	}

	if len(schedules) == 0 {
		return constant.CarrotScheduleShowEmpty
	}

	message := constant.CarrotScheduleShowSuccess
	for _, object := range schedules {
		duration := time.Now().Sub(object.Date.Local()).Hours() / 24
		message = fmt.Sprintf("%s\n[%s] %s %s %.0fd",
			message, object.ScheduleId, object.Date.Local().Format(timePatternSimple), object.Title, duration)
	} // [asd]1/2/Mon/数据结构考试

	return message
}

func getScheduleStringByTitleFromNow(title string, page int64, limit int64) string {
	m := model.GetModel()
	defer m.Close()

	schedules, err := m.GetScheduleByTitleFromNow(title, page, limit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return constant.CarrotScheduleNotFound
		}
		util.ErrorPrint(err, nil, "mongo")
		return constant.CarrotScheduleShowFailed
	}

	if len(schedules) == 0 {
		return constant.CarrotScheduleShowEmpty
	}

	message := constant.CarrotScheduleShowSuccess
	for _, object := range schedules {
		duration := time.Now().Sub(object.Date.Local()).Hours() / 24
		message = fmt.Sprintf("%s\n[%s] %s %s %.0fd",
			message, object.ScheduleId, object.Date.Local().Format(timePatternSimple), object.Title, duration)
	} // [asd]1/2/Mon/数据结构考试

	return message
}

func getScheduleDetailStringById(id string) string {
	m := model.GetModel()
	defer m.Close()

	object, err := m.GetScheduleById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return constant.CarrotScheduleNotFound
		}
		util.ErrorPrint(err, nil, "mongo")
		return constant.CarrotScheduleShowFailed
	}

	duration := time.Now().Sub(object.Date.Local())
	message := fmt.Sprintf("【第 %s 号约定】\n时间：%s\n剩余：%.1fh, %.0fd\n事件：%s\n备注：%s",
		object.ScheduleId, object.Date.Local().Format(timePatternComplex), duration.Hours(), duration.Hours()/24, object.Title, object.Description)

	return message
}

func getScheduleStringRecent() string {
	m := model.GetModel()
	defer m.Close()

	schedules, err := m.GetScheduleAllFromNow()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return constant.CarrotScheduleNotFound
		}
		util.ErrorPrint(err, nil, "mongo")
		return constant.CarrotScheduleShowFailed
	}

	if len(schedules) == 0 {
		return constant.CarrotScheduleShowEmpty
	}

	message := "哇！距离约定"
	flag := 0
	for _, object := range schedules {
		duration := time.Now().Sub(object.Date.Local()).Hours()
		if math.Abs(duration) < 24 && object.Alarm24hFlag == false && flag != 1 {
			message = fmt.Sprintf("%s「%s」", message, object.Title)
			object.Alarm24hFlag = true
			_, _ = m.UpdateScheduleById(object.ScheduleId, object)
			flag = 24
		} else if math.Abs(duration) < 1 && object.Alarm1hFlag == false && flag != 24 {
			message = fmt.Sprintf("%s「%s」", message, object.Title)
			object.Alarm1hFlag = true
			_, _ = m.UpdateScheduleById(object.ScheduleId, object)
			flag = 1
		}
	}
	message = fmt.Sprintf("%s只有不到 %d 小时了！卡洛好期待~", message, flag)

	return message
}
