package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func bootConfig() ConfigFile {
	var path = Utils.GetExecuteDir()
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
