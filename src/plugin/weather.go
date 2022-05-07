package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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
	if config.C.RiskControl {
		return false
	}
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
	util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getWeatherInfoString(location))
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
			FlagCanMatchedGroup:   !config.C.RiskControl,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         false,
		},
	}
	controller.PluginRegister(p)
}

func getWeatherInfoString(location string) string {
	url := config.C.Plugin.Weather.Host
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))

	q := req.URL.Query()
	q.Add("key", config.C.Plugin.Weather.Token)
	q.Add("location", location)
	q.Add("language", "zh-Hans")
	q.Add("unit", "c")
	q.Add("start", "0")
	q.Add("days", "3")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.ErrorPrint(err, nil, "weather get")
		return constant.CarrotWeatherFailed
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.ErrorPrint(err, nil, "weather get")
		return constant.CarrotWeatherFailed
	}

	var responseWeather param.ResponseWeather
	if err = json.Unmarshal(body, &responseWeather); err != nil {
		util.ErrorPrint(err, nil, "weather get")
		return constant.CarrotWeatherFailed
	}

	if len(responseWeather.Result) == 0 {
		util.ErrorPrint(errors.New("zero"), nil, "weather get")
		return constant.CarrotWeatherFailed
	}

	message := constant.CarrotWeatherStart
	for _, wt := range responseWeather.Result[0].Daily {
		val := reflect.ValueOf(wt)
		typ := reflect.TypeOf(wt)
		num := val.NumField()

		message += fmt.Sprintf("\n%s %s 天气: ", responseWeather.Result[0].Location.Name, wt.Date)
		for i := 0; i < num; i++ {
			if typ.Field(i).Tag.Get("dismiss") != "true" && val.Field(i).String() != "" {
				message += fmt.Sprintf("%s「%s」", typ.Field(i).Tag.Get("text"), strings.Replace(val.Field(i).String(), "\n", "", -1))
			}
		}
	}

	return message
}
