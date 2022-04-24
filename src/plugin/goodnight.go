package plugin

import (
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
	"time"
)

type goodNight struct {
	Index       param.PluginIndex
	PassHour    map[int]bool
	UserDay     map[int64]int
	LastUserDay int
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

	// not night
	if !exist || !ok {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNightCheat)
		return nil
	}

	// already greeting
	userDay, exist := p.UserDay[util.GetQQGroupUserId(msg)]
	if exist && userDay == day {
		util.QQGroupSendAtSomeone(msg.GroupId, id, constant.CarrotGroupGoodNightRepeat)
		return nil
	}

	//
	p.UserDay[id] = day
	if p.LastUserDay != day {
		p.LastUserDay = day
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
			FlagCanListen:         false,
		},
		PassHour:    passHour,
		UserDay:     make(map[int64]int),
		LastUserDay: time.Now().Day() - 1,
	}
	controller.PluginRegister(p)
}
