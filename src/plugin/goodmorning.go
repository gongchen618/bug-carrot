package plugin

import (
	"bug-carrot/src/config"
	"bug-carrot/src/constant"
	"bug-carrot/src/controller"
	"bug-carrot/src/param"
	"bug-carrot/src/util"
	"time"
)

type goodMorning struct {
	Index                  param.PluginIndex
	PassHour               map[int]bool
	UserDay                map[int64]int
	LastUserDay            int
	LastAutoGoodMorningDay int
}

func (p *goodMorning) GetPluginName() string {
	return p.Index.PluginName
}
func (p *goodMorning) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *goodMorning) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *goodMorning) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *goodMorning) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *goodMorning) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *goodMorning) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *goodMorning) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

func (p *goodMorning) IsTime() bool {
	if time.Now().Hour() == 7 && time.Now().Day() != p.LastAutoGoodMorningDay {
		p.LastAutoGoodMorningDay = time.Now().Day()
		return true
	}
	return false
}
func (p *goodMorning) DoTime() error {
	util.QQGroupSend(config.C.Plugin.Schedule.Group, "卡洛起床啦！大家早安！")
	return nil
}

func (p *goodMorning) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"早安"})
}
func (p *goodMorning) DoMatchedGroup(msg param.GroupMessage) error {
	hour, day := time.Now().Hour(), time.Now().Day()
	ok, exist := p.PassHour[hour]
	id := util.GetQQGroupUserId(msg)

	// not morning
	if !exist || !ok {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodMorningCheat)
		return nil
	}

	// already greeting
	userDay, exist := p.UserDay[util.GetQQGroupUserId(msg)]
	if exist && userDay == day {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodMorningRepeat)
		return nil
	}

	//
	p.UserDay[id] = day
	if p.LastUserDay != day {
		p.LastUserDay = day
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodMorningFirst)
	} else {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodMorning)
	}
	return nil
}

func (p *goodMorning) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return config.C.RiskControl && msg.WordsMap.ExistWord("n", []string{"早安"})
}
func (p *goodMorning) DoMatchedPrivate(msg param.PrivateMessage) error {
	hour, day := time.Now().Hour(), time.Now().Day()
	ok, exist := p.PassHour[hour]
	id := msg.UserId

	// not morning
	if !exist || !ok {
		util.QQSend(id, constant.CarrotGroupGoodMorningCheat)
		return nil
	}

	// already greeting
	userDay, exist := p.UserDay[id]
	if exist && userDay == day {
		util.QQSend(id, constant.CarrotGroupGoodMorningRepeat)
		return nil
	}

	//
	p.UserDay[id] = day
	if p.LastUserDay != day {
		p.LastUserDay = day
		util.QQSend(id, constant.CarrotGroupGoodMorningFirst)
	} else {
		util.QQSend(id, constant.CarrotGroupGoodMorning)
	}

	return nil
}

func (p *goodMorning) Listen(req param.GroupMessage) {
}

func (p *goodMorning) Close() {
}

func GoodMorningPluginRegister() {
	passHour := make(map[int]bool)
	hour := []int{5, 6, 7, 8, 9}
	for _, h := range hour {
		passHour[h] = true
	}
	p := &goodMorning{
		Index: param.PluginIndex{
			PluginName:            "goodMorning",
			PluginAuthor:          "gongchen618",
			FlagCanTime:           true,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         false,
			FlagUseDatabase:       false,
			FlagIgnoreRiskControl: false,
		},
		PassHour:               passHour,
		UserDay:                make(map[int64]int),
		LastUserDay:            time.Now().Day() - 1,
		LastAutoGoodMorningDay: time.Now().Day() - 1,
	}
	controller.PluginRegister(p)
}
