package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConfigFile struct {
	Debug    bool   `json:"Debug"`
	Timezone string `json:"Timezone"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DB       struct {
		Driver   string `json:"driver"`
		Postgres struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			DbName   string `json:"database"`
		} `json:"postgres"`
	} `json:"DB"`
	Logger struct {
		WriteToConsole bool `json:"writeToConsole"`
		WriteToFile    struct {
			Active      bool  `json:"active"`
			MaxLogFiles int   `json:"maxLogFiles"`
			MaxLogSize  int64 `json:"maxLogSize"`
		} `json:"writeToFile"`
	} `json:"Logger"`
}

func bootConfig() ConfigFile {
	var path = servus.Utils.GetExecuteDir()
	var config ConfigFile
	path = fmt.Sprintf("%v/config.json", path)
	byteValue, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}
	return config
}
