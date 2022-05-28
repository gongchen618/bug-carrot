package plugin

import (
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

type dice struct {
	Index          param.PluginIndex
	DicePrefix     string
	DividingString string
}

func (p *dice) GetPluginName() string {
	return p.Index.PluginName
}
func (p *dice) GetPluginAuthor() string {
	return p.Index.PluginAuthor
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
func (p *dice) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *dice) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

func (p *dice) IsTime() bool {
	return false
}
func (p *dice) DoTime() error {
	return nil
}

func (p *dice) IsMatchedGroup(msg param.GroupMessage) bool { // 占卜[name]
	if strings.HasPrefix(msg.RawMessage, p.DicePrefix) {
		return true
	}
	return false
}
func (p *dice) DoMatchedGroup(msg param.GroupMessage) error {
	topic := msg.RawMessage[len(p.DicePrefix):]
	limitTag := false
	limit := int64(0)
	var err error
	if strings.Count(topic, p.DividingString) == 1 {
		str := strings.Split(topic, p.DividingString)
		topic = str[0]
		limit, err = strconv.ParseInt(str[1], 10, 64)
		if err != nil || limit <= 0 {
			util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), "好奇怪的要求...")
			return nil
		} else {
			limitTag = true
		}
	} // 占卜topic#number

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

	if limitTag { // 有上界的占卜
		rd, err := rand.Int(rand.Reader, big.NewInt(limit))
		if err != nil { // 在最上面检验了 limit 是否为正数，所以这个 err 应该恒为 nil
			rd = big.NewInt(0)
		}
		star := rd.Int64() + 1 // 下界从 1 开始
		diceResultMessage := fmt.Sprintf("#卡洛透过水晶球向 %d 颗星星望去，发现与事件「%s」最契合的是小行星 %d 号...这意味着什么呢？", limit, topic, star)
		util.QQGroupSend(msg.GroupId, diceResultMessage)

		return nil
	}
	fmt.Println(topic)

	rd, err := rand.Int(rand.Reader, big.NewInt(101))
	if err != nil {
		rd = big.NewInt(0)
	}
	star := rd.Int64()
	var rdLevelMessage string
	switch {
	case star == 100:
		rdLevelMessage = constant.CarrotDiceSuccessFullPoint
	case star >= 95:
		rdLevelMessage = constant.CarrotDiceSuccessGold
	case star >= 85:
		rdLevelMessage = constant.CarrotDiceSuccessSilver
	case star >= 60:
		rdLevelMessage = constant.CarrotDiceSuccessBronze
	case star >= 40:
		rdLevelMessage = constant.CarrotDiceFailedGold
	case star > 0:
		rdLevelMessage = constant.CarrotDiceFailedSilver
	case star == 0:
		rdLevelMessage = constant.CarrotDiceFailedZeroPoint
	}

	diceResultMessage := fmt.Sprintf("#卡洛对事件「%s」使用了占卜术，一共有 %d 颗星星被点亮，星象显示「%s」", topic, star, rdLevelMessage)
	util.QQGroupSend(msg.GroupId, diceResultMessage)

	return nil
}

func (p *dice) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return false
}
func (p *dice) DoMatchedPrivate(msg param.PrivateMessage) error {
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
			PluginAuthor:          "gongchen618",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         false,
			FlagUseDatabase:       false,
			FlagIgnoreRiskControl: false,
		},
		DicePrefix:     "占卜",
		DividingString: "#",
	}
	controller.PluginRegister(p)
}
