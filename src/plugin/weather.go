package plugin

import (
	"bug-carrot/controller"
	"bug-carrot/controller/param"
	"bug-carrot/util"
)

type weather struct {
	PluginName string
}

func (p *weather) GetPluginName() string {
	return p.PluginName
}

func (p *weather) IsTime() bool {
	return false
}

func (p *weather) DoTime() error {
	return nil
}

func (p *weather) IsMatched(msg param.GroupMessage) bool {
	return util.IsWordInMessage("n", []string{"天气"}, msg)
}

func (p *weather) DoMatched(msg param.GroupMessage) error {
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

func (p *weather) Listen(msg param.GroupMessage) {

}

func (p *weather) Close() {
}

func WeatherPluginRegister() {
	p := &weather{
		PluginName: "weather",
	}
	controller.PluginRegister(p)
}
