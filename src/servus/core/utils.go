package core

import "os"

type utils interface {
	GetExecuteDir() string
}

type Utils struct {
	utils
}

func (u Utils) GetExecuteDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}