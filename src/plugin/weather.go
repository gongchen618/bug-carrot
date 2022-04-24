package plugin

import (
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
)

type weather struct {
	Index param.PluginIndex
}

func (p *weather) GetPluginName() string {
	return p.Index.PluginName
}
func (p *weather) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *weather) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *weather) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *weather) CanListen() bool {
	return p.Index.FlagCanListen
}

func (p *weather) IsTime() bool {
	return false
}
func (p *weather) DoTime() error {
	return nil
}

func (p *weather) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.WordsMap.ExistWord("n", []string{"天气"})
}
func (p *weather) DoMatchedGroup(msg param.GroupMessage) error {
	location := "武汉"
	words := util.GetWordsFromMessage(msg.RawMessage)
	for _, word := range words {
		if word.Type == "ns" {
			location = word.Word
		}
	}
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), util.GetWeatherInfoString(location))
	return nil
}

func (p *weather) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return false
}
func (p *weather) DoMatchedPrivate(msg param.PrivateMessage) error {
	return nil
}

func (p *weather) Listen(msg param.GroupMessage) {
}

func (p *weather) Close() {
}

func WeatherPluginRegister() {
	p := &weather{
		Index: param.PluginIndex{
			PluginName:            "weather",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         false,
		},
	}
	controller.PluginRegister(p)
}
