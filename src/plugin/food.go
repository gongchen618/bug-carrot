package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/model"
	"bug-carrot/param"
	"bug-carrot/util"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type food struct {
	Index                param.PluginIndex
	FoodAddPrefix        string
	FoodDeletePrefix     string
	FoodQueryPrefix      string
	DividingString       string
	FoodAddPrefixPrivate string
}

func (p *food) GetPluginName() string {
	return p.Index.PluginName
}
func (p *food) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *food) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *food) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *food) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *food) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *food) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *food) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

func (p *food) IsTime() bool {
	return false
}
func (p *food) DoTime() error {
	return nil
}

func (p *food) IsMatchedGroup(msg param.GroupMessage) bool {
	if msg.Anonymous.Id != 0 { // 禁止匿名
		return false
	}
	if strings.HasPrefix(msg.RawMessage, p.FoodAddPrefix) ||
		strings.HasPrefix(msg.RawMessage, p.FoodDeletePrefix) ||
		strings.HasPrefix(msg.RawMessage, p.FoodQueryPrefix) {
		return true
	}
	return false
}
func (p *food) DoMatchedGroup(msg param.GroupMessage) error {
	switch {
	case strings.HasPrefix(msg.RawMessage, p.FoodAddPrefix): // 格式：安利name#place#description
		info := msg.RawMessage[len(p.FoodAddPrefix):]
		if strings.Count(info, p.DividingString) != 2 {
			util.QQGroupSendAtSomeone(msg.GroupId, msg.UserId, constant.CarrotFoodStrangeInput)
			return nil
		}
		str := strings.Split(info, p.DividingString)
		foodAddGroup(msg.UserId, msg.GroupId, str[0], str[1], str[2])

	case strings.HasPrefix(msg.RawMessage, p.FoodDeletePrefix): // 格式：拔草name
		name := msg.RawMessage[len(p.FoodDeletePrefix):]
		foodDelete(msg.GroupId, name)

	case strings.HasPrefix(msg.RawMessage, p.FoodQueryPrefix): // 格式：吃什么 / 吃什么#place
		info := msg.RawMessage[len(p.FoodQueryPrefix):]
		if strings.Count(info, p.DividingString) == 1 {
			str := strings.Split(msg.RawMessage, p.DividingString)
			foodRandByPlaceGroup(msg.UserId, msg.GroupId, str[1])
			return nil
		}
		foodRandAllGroup(msg.UserId, msg.GroupId)
	}
	return nil
}

func (p *food) IsMatchedPrivate(msg param.PrivateMessage) bool {
	if msg.UserId == config.C.Plugin.Food.Admin && strings.HasPrefix(msg.RawMessage, "查杀") {
		return true
	}
	if strings.HasPrefix(msg.RawMessage, p.FoodAddPrefixPrivate) {
		return true
	}
	if config.C.RiskControl && strings.HasPrefix(msg.RawMessage, p.FoodQueryPrefix[1:]) {
		return true
	}
	return false
}
func (p *food) DoMatchedPrivate(msg param.PrivateMessage) error { // 格式：查杀 xx
	if msg.UserId == config.C.Plugin.Food.Admin && strings.HasPrefix(msg.RawMessage, "查杀") {
		str := strings.Split(msg.RawMessage, " ") // 没有考虑错误情况 因为是 admin private message
		if len(str) >= 2 {
			foodCheck(msg.UserId, str[1])
			return nil
		}
		util.QQSend(msg.UserId, constant.CarrotGroupPuzzled)
	} else if strings.HasPrefix(msg.RawMessage, p.FoodAddPrefixPrivate) {
		info := msg.RawMessage[len(p.FoodAddPrefixPrivate):]
		if strings.Count(info, p.DividingString) != 2 {
			util.QQSend(msg.UserId, constant.CarrotFoodStrangeInput)
			return nil
		}
		str := strings.Split(info, p.DividingString)
		foodAddPrivate(msg.UserId, str[0], str[1], str[2])
	} else {
		info := msg.RawMessage[len(p.FoodQueryPrefix[1:]):]
		if strings.Count(info, p.DividingString) == 1 {
			str := strings.Split(msg.RawMessage, p.DividingString)
			foodRandByPlacePrivate(msg.UserId, str[1])
			return nil
		}
		foodRandAllGroupPrivate(msg.UserId)
	}
	return nil
}

func (p *food) Listen(msg param.GroupMessage) {

}

func (p *food) Close() {
}

func FoodPluginRegister() {
	p := &food{
		Index: param.PluginIndex{
			PluginName:            "food",
			PluginAuthor:          "gongchen618",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         false,
			FlagUseDatabase:       true,
			FlagIgnoreRiskControl: false,
		},
		FoodAddPrefix:        "安利",
		FoodDeletePrefix:     "拔草",
		FoodQueryPrefix:      "吃什么",
		FoodAddPrefixPrivate: "安利",
		DividingString:       "#",
	}
	controller.PluginRegister(p)
}

func foodAddGroup(id int64, group int64, name string, address string, description string) {
	m := model.GetModel()
	defer m.Close()

	fd := param.Food{
		Name:        name,
		Address:     address,
		Description: description,
		Recommender: id,
	}

	if err := m.AddFood(fd); err != nil {
		util.QQGroupSend(group, constant.CarrotFoodAddFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQGroupSend(group, fmt.Sprintf("卡洛接受了「%s」的安利！有机会一起去尝尝看吧！", name))
}

func foodAddPrivate(id int64, name string, address string, description string) {
	m := model.GetModel()
	defer m.Close()

	fd := param.Food{
		Name:        name,
		Address:     address,
		Description: description,
		Recommender: id,
	}

	if err := m.AddFood(fd); err != nil {
		util.QQSend(id, constant.CarrotFoodAddFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, fmt.Sprintf("卡洛接受了「%s」的安利！有机会一起去尝尝看吧！", name))
}

func foodDelete(group int64, name string) {
	m := model.GetModel()
	defer m.Close()

	fd := param.Food{
		Name: name,
	}

	if err := m.DeleteFood(fd); err != nil {
		util.QQGroupSend(group, constant.CarrotFoodDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQGroupSend(group, fmt.Sprintf("「%s」被拔草啦！", name))
}

func foodRandByPlaceGroup(id int64, group int64, place string) {
	util.QQGroupSendAtSomeone(group, id, getRandomFoodStringByPlace(place))
}

func foodRandByPlacePrivate(id int64, place string) {
	util.QQSend(id, getRandomFoodStringByPlace(place))
}

func foodRandAllGroup(id int64, group int64) {
	util.QQGroupSendAtSomeone(group, id, getRandomFoodString())
}

func foodRandAllGroupPrivate(id int64) {
	util.QQSend(id, getRandomFoodString())
}

func foodCheck(id int64, name string) {
	util.QQSend(id, getFoodInformation(name))
}

func getRandomFoodStringByPlace(place string) string {
	m := model.GetModel()
	defer m.Close()

	foods, err := m.GetFoodByAddress(place)
	if err != nil {
		util.ErrorPrint(err, place, "mongo")
		return constant.CarrotFoodRandFailed
	}
	if len(foods) == 0 {
		return constant.CarrotFoodRandEmpty
	}

	rand.Seed(time.Now().UnixNano())
	rd := rand.Intn(len(foods))
	message := fmt.Sprintf("想在「%s」吃的话，卡洛推荐「%s」哦！它「%s」！", place, foods[rd].Name, foods[rd].Description)

	return message
}

func getRandomFoodString() string {
	m := model.GetModel()
	defer m.Close()

	foods, err := m.GetFoodAll()
	if err != nil {
		util.ErrorPrint(err, nil, "mongo")
		return constant.CarrotFoodRandFailed
	}
	if len(foods) == 0 {
		return constant.CarrotFoodRandEmpty
	}

	rand.Seed(time.Now().UnixNano())
	rd := rand.Intn(len(foods))
	message := fmt.Sprintf("饿了的话，试试「%s」的「%s」怎么样？它「%s」~", foods[rd].Address, foods[rd].Name, foods[rd].Description)

	return message
}

func getFoodInformation(name string) string {
	m := model.GetModel()
	defer m.Close()

	foods, err := m.GetFoodByName(name)
	if err != nil {
		util.ErrorPrint(err, name, "mongo")
		return constant.CarrotFoodInfoGetFailed
	}
	if len(foods) == 0 {
		return constant.CarrotFoodInfoGetEmpty
	}

	message := fmt.Sprintf("查询结果：名称「%s」，地址「%s」，描述「%s」，推荐人「%d」",
		foods[0].Name, foods[0].Address, foods[0].Description, foods[0].Recommender)

	return message
}
