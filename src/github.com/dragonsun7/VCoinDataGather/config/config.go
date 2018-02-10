package config

import (
	"sync"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	DB struct {
		Driver   string `json:"driver"`
		Host     string `json:"host"`
		Port     int64  `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	}
}

const (
	kConfigFile = "./config.json"
)

var (
	sharedInstance *Config
	once           sync.Once
)

// 单例
func GetInstance() (*Config) {
	once.Do(func() {
		sharedInstance = new(Config)

		data, err := ioutil.ReadFile(kConfigFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(data, sharedInstance)
		if err != nil {
			panic(err)
		}
	})

	return sharedInstance
}
