package _core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConfigFile struct {
	Debug    bool   `json:"Debug"`
	Timezone string `json:"Timezone"`
	Host string `json:"Host"`
	Port string `json:"Port"`
	DB       struct {
		Driver   string `json:"Driver"`
		Postgres struct {
			Host     string `json:"Host"`
			Port     string `json:"Port"`
			User     string `json:"User"`
			Password string `json:"Password"`
			DbName   string `json:"Database"`
		} `json:"Postgres"`
	} `json:"DB"`
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
