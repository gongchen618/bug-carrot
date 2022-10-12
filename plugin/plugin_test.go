package my_plugin

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
	"plugin"
	"testing"
)

type PluginInfo struct {
	PluginName string `yaml:"plugin-name"`
	Method     method `yaml:"method"`
}

type method struct {
	FunctionName string `yaml:"function-name"`
}

func getConfig() *PluginInfo {
	p, err := os.Getwd()
	if err != nil {
		log.Println("ERR", err.Error())
	}
	pluginFilePath := path.Join(p, "dice/dice.yml")
	data, err := ioutil.ReadFile(pluginFilePath)
	if err != nil {
		log.Println("ERR", err.Error())
	}

	pluginInfo := &PluginInfo{}
	err = yaml.Unmarshal(data, pluginInfo)
	if err != nil {
		log.Println("ERR", err.Error())
	}
	fmt.Println(pluginInfo)
	return pluginInfo
}

func TestDice(t *testing.T) {
	pluginInfo := getConfig()

	plu, err := plugin.Open("plugin-dice.so")
	if err != nil {
		log.Println("ERR", err.Error())
	}

	pluFunc, err := plu.Lookup(pluginInfo.Method.FunctionName)
	if err != nil {
		log.Println("ERR", err.Error())
	}

	pluFunc.(func(string))("234")
}
