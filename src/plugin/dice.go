package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

type dice struct {
	Index      param.PluginIndex
	DicePrefix string
}

func (p *dice) GetPluginName() string {
	return p.Index.PluginName
}
func (p *dice) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *dice) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *dice) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *dice) CanListen() bool {
	return p.Index.FlagCanListen
}

func (p *dice) IsTime() bool {
	return false
}
func (p *dice) DoTime() error {
	return nil
}

func (p *dice) IsMatchedGroup(msg param.GroupMessage) bool { // 占卜[name]
	if config.C.RiskControl {
		return false
	}
	if strings.HasPrefix(msg.RawMessage, p.DicePrefix) {
		return true
	}
	return false
}
func (p *dice) DoMatchedGroup(msg param.GroupMessage) error {
	topic := msg.RawMessage[len(p.DicePrefix):]
	for _, ch := range topic {
		if !unicode.Is(unicode.Han, ch) && !unicode.IsLetter(ch) && !unicode.IsNumber(ch) {
			util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotDiceStrangeInput)
			return nil
		}
	}

	if len(topic) == 0 {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotDiceEmptyTopic)
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	rd := rand.Intn(101)
	var rdLevelMessage string
	switch {
	case rd == 100:
		rdLevelMessage = constant.CarrotDiceSuccessFullPoint
	case rd >= 95:
		rdLevelMessage = constant.CarrotDiceSuccessGold
	case rd >= 85:
		rdLevelMessage = constant.CarrotDiceSuccessSilver
	case rd >= 60:
		rdLevelMessage = constant.CarrotDiceSuccessBronze
	case rd >= 40:
		rdLevelMessage = constant.CarrotDiceFailedGold
	case rd > 0:
		rdLevelMessage = constant.CarrotDiceFailedSilver
	case rd == 0:
		rdLevelMessage = constant.CarrotDiceFailedZeroPoint
	}

	diceResultMessage := fmt.Sprintf("#卡洛对事件「%s」使用了占卜术，一共有 %d 颗星星被点亮，星象显示「%s」", topic, rd, rdLevelMessage)
	util.QQGroupSend(msg.GroupId, diceResultMessage)

	return nil
}

func (p *dice) IsMatchedPrivate(msg param.PrivateMessage) bool {
	if config.C.RiskControl && strings.HasPrefix(msg.RawMessage, p.DicePrefix[1:]) {
		return true
	}
	return false
}
func (p *dice) DoMatchedPrivate(msg param.PrivateMessage) error {
	topic := msg.RawMessage[len(p.DicePrefix[1:]):]
	for _, ch := range topic {
		if !unicode.Is(unicode.Han, ch) && !unicode.IsLetter(ch) && !unicode.IsNumber(ch) {
			util.QQSend(msg.UserId, constant.CarrotDiceStrangeInput)
			return nil
		}
	}

	if len(topic) == 0 {
		util.QQSend(msg.UserId, constant.CarrotDiceEmptyTopic)
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	rd := rand.Intn(101)
	var rdLevelMessage string
	switch {
	case rd == 100:
		rdLevelMessage = constant.CarrotDiceSuccessFullPoint
	case rd >= 95:
		rdLevelMessage = constant.CarrotDiceSuccessGold
	case rd >= 85:
		rdLevelMessage = constant.CarrotDiceSuccessSilver
	case rd >= 60:
		rdLevelMessage = constant.CarrotDiceSuccessBronze
	case rd >= 40:
		rdLevelMessage = constant.CarrotDiceFailedGold
	case rd > 0:
		rdLevelMessage = constant.CarrotDiceFailedSilver
	case rd == 0:
		rdLevelMessage = constant.CarrotDiceFailedZeroPoint
	}

	diceResultMessage := fmt.Sprintf("#卡洛对事件「%s」使用了占卜术，一共有 %d 颗星星被点亮，星象显示「%s」", topic, rd, rdLevelMessage)
	util.QQSend(msg.UserId, diceResultMessage)

	return nil
}

func (p *dice) Listen(msg param.GroupMessage) {

}

func (p *dice) Close() {
}

func DicePluginRegister() {
	p := &dice{
		Index: param.PluginIndex{
			PluginName:            "dice",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   !config.C.RiskControl,
			FlagCanMatchedPrivate: config.C.RiskControl,
			FlagCanListen:         false,
		},
		DicePrefix: " 占卜",
	}
	controller.PluginRegister(p)
}
