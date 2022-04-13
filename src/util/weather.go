package util

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller/param"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func GetWeatherInfoString(location string) string {
	url := fmt.Sprintf("%s", config.C.Weather.Host)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))

	q := req.URL.Query()
	q.Add("key", config.C.Weather.Token)
	q.Add("location", location)
	q.Add("language", "zh-Hans")
	q.Add("unit", "c")
	q.Add("start", "0")
	q.Add("days", "3")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ErrorPrint(err, nil, "weather get")
		return constant.CarrotWeatherFailed
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorPrint(err, nil, "weather get")
		return constant.CarrotWeatherFailed
	}

	var responseWeather param.ResponseWeather
	if err = json.Unmarshal(body, &responseWeather); err != nil {
		ErrorPrint(err, nil, "weather get")
		return constant.CarrotWeatherFailed
	}

	if len(responseWeather.Result) == 0 {
		ErrorPrint(errors.New("zero"), nil, "weather get")
		return constant.CarrotWeatherFailed
	}

	message := constant.CarrotWeatherStart
	for _, weather := range responseWeather.Result[0].Daily {
		val := reflect.ValueOf(weather)
		typ := reflect.TypeOf(weather)
		num := val.NumField()

		message += fmt.Sprintf("\n%s %s 天气: ", responseWeather.Result[0].Location.Name, weather.Date)
		for i := 0; i < num; i++ {
			if typ.Field(i).Tag.Get("dismiss") != "true" && val.Field(i).String() != "" {
				message += fmt.Sprintf("%s「%s」", typ.Field(i).Tag.Get("text"), strings.Replace(val.Field(i).String(), "\n", "", -1))
			}
		}
	}

	return message
}
