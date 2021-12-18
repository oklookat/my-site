package core

import (
	"io"
	"servus/core/internal/controlify/controlifyTelegram"
)

type controlTelegram struct {
	config *ControlTelegramConfig
	logger Logger
	Bot    controlifyTelegram.Controller
}

func (c *controlTelegram) new(config *ControlTelegramConfig, logger Logger) {
	if config == nil {
		panic("[control/telegram]: config nil pointer.")
	}
	if logger == nil {
		panic("[control/telegram]: logger nil pointer.")
	}
	c.config = config
	c.logger = logger
	c.Bot = controlifyTelegram.Controller{}
	c.Bot.New(c)
}

func (c *controlTelegram) checkErr(err error) {
	if err != nil {
		c.logger.Error("[control/telegram]: " + err.Error())
	}
}

func (c *controlTelegram) SendFile(caption *string, filename string, reader io.Reader) {
	err := c.Bot.SendFile(caption, filename, reader)
	c.checkErr(err)
}

func (c *controlTelegram) SendMessage(message string) {
	err := c.Bot.SendMessage(message)
	c.checkErr(err)
}

func (c *controlTelegram) GetToken() string {
	return c.config.Token
}

func (c *controlTelegram) GetAllowedChats() []int64 {
	return c.config.AllowedChats
}

func (c *controlTelegram) GetAllowedUsers() []int64 {
	return c.config.AllowedUsers
}
