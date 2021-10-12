package core

import (
	"encoding/json"
	"fmt"
	"os"
)

func bootConfig() *ConfigFile {
	var path = Utils.GetExecuteDir()
	var config ConfigFile
	path = fmt.Sprintf("%v/settings/config.json", path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
