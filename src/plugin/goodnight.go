package plugin

import (
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
	"fmt"
	"time"
)

type goodNight struct {
	Index                param.PluginIndex
	PassHour             map[int]bool
	UserLastGoodNightDay map[int64]int
	UserMessageCount     map[int64]int
	LastGoodNightDay     int
	TimeDividingLine     int
}

func (p *goodNight) GetPluginName() string {
	return p.Index.PluginName
}
func (p *goodNight) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *goodNight) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *goodNight) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *goodNight) CanListen() bool {
	return p.Index.FlagCanListen
}

func (p *goodNight) IsTime() bool {
	return false
}
func (p *goodNight) DoTime() error {
	return nil
}

func (p *goodNight) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"晚安"})
}
func (p *goodNight) DoMatchedGroup(msg param.GroupMessage) error {
	hour, day := time.Now().Hour(), time.Now().Day()
	ok, exist := p.PassHour[hour]
	id := util.GetQQGroupUserId(msg)

	// 处理 12 点之后是昨天
	if hour <= p.TimeDividingLine {
		d, _ := time.ParseDuration("-24h")
		day = time.Now().Add(d).Day()
	}

	// not night
	if !exist || !ok {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNightCheat)
		return nil
	}

	// already greeting
	userDay, exist := p.UserLastGoodNightDay[util.GetQQGroupUserId(msg)]
	if exist && userDay == day {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNightRepeat)
		return nil
	}

	p.UserLastGoodNightDay[id] = day
	if p.LastGoodNightDay != day {
		p.LastGoodNightDay = day
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNightFirst)
	} else {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNight)
	}
	return nil
}

func (p *goodNight) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return false
}
func (p *goodNight) DoMatchedPrivate(msg param.PrivateMessage) error {
	return nil
}

func (p *goodNight) Listen(msg param.GroupMessage) {
	hour, day := time.Now().Hour(), time.Now().Day()
	ok, exist := p.PassHour[hour]
	id := util.GetQQGroupUserId(msg)

	// 处理 12 点之后是昨天
	if hour <= p.TimeDividingLine {
		d, _ := time.ParseDuration("-24h")
		day = time.Now().Add(d).Day()
	}

	// not night
	if !exist || !ok {
		return
	}

	// already greeting but chat
	userDay, exist := p.UserLastGoodNightDay[util.GetQQGroupUserId(msg)]
	if exist && userDay == day {
		cnt, flag := p.UserMessageCount[id]
		if flag {
			p.UserMessageCount[id] = cnt + 1
		} else {
			p.UserMessageCount[id] = 1
		}

		if p.UserMessageCount[id]%5 == 0 {
			util.QQGroupBan(msg.GroupId, id, int64(60*(p.UserMessageCount[id]/5)))
			util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNightButChat)
		} else if p.UserMessageCount[id]%4 == 0 {
			util.QQGroupSend(msg.GroupId, fmt.Sprintf("[CQ:poke,qq=%d]", id))
		}
		return
	}
}

func (p *goodNight) Close() {
}

func GoodNightPluginRegister() {
	passHour := make(map[int]bool)
	hour := []int{22, 23, 24, 0, 1, 2, 3, 4, 5}
	for _, h := range hour {
		passHour[h] = true
	}
	p := &goodNight{
		Index: param.PluginIndex{
			PluginName:            "goodnight",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         true,
		},
		PassHour:             passHour,
		UserLastGoodNightDay: make(map[int64]int),
		UserMessageCount:     make(map[int64]int),
		LastGoodNightDay:     time.Now().Day() - 1,
		TimeDividingLine:     6,
	}
	controller.PluginRegister(p)
}
