package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	C *Config
)

type Config struct {
	App         app     `yaml:"app"`
	MongoDB     mongodb `yaml:"mongodb"`
	QQBot       qqbot   `yaml:"qqbot"`
	Plugin      plugin  `yaml:"plugin"`
	RiskControl bool    `yaml:"risk-control"`
	DatabaseUse bool    `yaml:"database-use"`
}

type mongodb struct {
	Host string `yaml:"host"`
}

type app struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}

type qqbot struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	QQ   int64  `yaml:"qq"`
}

type plugin struct {
	Default    _default   `yaml:"default"`
	Weather    weather    `yaml:"weather"`
	Homework   homework   `yaml:"homework"`
	Food       food       `yaml:"food"`
	Schedule   schedule   `yaml:"schedule"`
	Codeforces codeforces `yaml:"codeforces"`
}

type _default struct {
	Admin int64 `yaml:"admin"`
}

type weather struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type homework struct {
	Admin int64 `yaml:"admin"`
	Group int64 `yaml:"group"`
}

type food struct {
	Admin int64 `yaml:"admin"`
	Group int64 `yaml:"group"`
}

type schedule struct {
	Admin int64 `yaml:"admin"`
	Group int64 `yaml:"group"`
}

type codeforces struct {
	Admin int64 `yaml:"admin"`
	Group int64 `yaml:"group"`
}

func init() {
	configFile := "config/default.yml"

	// 如果设置了
	if v, ok := os.LookupEnv("CONFIG"); ok {
		configFile = "config/" + v + ".yml"
	}

	configFilePath := ""
	p, err := os.Getwd()
	if err != nil {
		log.Panic(err)
		return
	}

	for {
		configFilePath = path.Join(p, configFile)
		if _, err := os.Stat(configFilePath); err == nil {
			break
		}
		if p == "/" {
			log.Panic("config not found")
		}
		p = path.Join(p, "..")
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("Read config error!")
		log.Panic(err)
		return
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	C = config

	log.Println("Config " + configFile + " loaded.")
}
