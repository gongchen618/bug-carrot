package dice

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

type PluginInfo struct {
	PluginName string `yaml:"plugin-name"`
	Method     method `yaml:"method"`
}

type method struct {
	FunctionName string `yaml:"function-name"`
}

func TestDice(t *testing.T) {
	p, err := os.Getwd()
	if err != nil {
		log.Println("ERR", err.Error())
	}
	pluginFilePath := path.Join(p, "dice.yml")
	data, err := ioutil.ReadFile(pluginFilePath)
	if err != nil {
		log.Println("ERR", err.Error())
	}
	plugin := &PluginInfo{}
	err = yaml.Unmarshal(data, plugin)
	if err != nil {
		log.Println("ERR", err.Error())
	}
	fmt.Println(plugin)

	dp := CallPlugin()
	value := reflect.ValueOf(&dp)
	f := value.MethodByName(plugin.Method.FunctionName) //通过反射获取它对应的函数，然后通过call来调用
	f.Call([]reflect.Value{reflect.ValueOf("123")})
}
