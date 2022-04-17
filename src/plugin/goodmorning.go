package plugin

import (
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/controller/param"
	"bug-carrot/util"
	"time"
)

type goodMorning struct {
	PluginName  string
	PassHour    map[int]bool
	UserDay     map[int64]int
	LastUserDay int
}

func (p *goodMorning) GetPluginName() string {
	return p.PluginName
}

func (p *goodMorning) IsTime() bool {
	return false
}

func (p *goodMorning) DoTime() error {
	return nil
}

func (p *goodMorning) IsMatched(req param.GroupMessage) bool {
	return util.IsWordInMessage("n", []string{"早安"}, req)
}

func (p *goodMorning) DoMatched(req param.GroupMessage) error {
	hour, day := time.Now().Hour(), time.Now().Day()
	ok, exist := p.PassHour[hour]
	id := util.GetQQGroupUserId(req)

	// not night
	if !exist || !ok {
		util.QQGroupSendAtSomeone(req.GroupId, id, constant.CarrotGroupGoodMorningCheat)
		return nil
	}

	// already greeting
	userDay, exist := p.UserDay[util.GetQQGroupUserId(req)]
	if exist && userDay == day {
		util.QQGroupSendAtSomeone(req.GroupId, id, constant.CarrotGroupGoodMorningRepeat)
		return nil
	}

	//
	p.UserDay[id] = day
	if p.LastUserDay != day {
		p.LastUserDay = day
		util.QQGroupSendAtSomeone(req.GroupId, id, constant.CarrotGroupGoodMorningFirst)
	} else {
		util.QQGroupSendAtSomeone(req.GroupId, id, constant.CarrotGroupGoodMorning)
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
		PluginName:  "goodMorning",
		PassHour:    passHour,
		UserDay:     make(map[int64]int),
		LastUserDay: time.Now().Day() - 1,
	}
	controller.PluginRegister(p)
}
