package core

import (
	"io"
	"servus/core/internal/controlify/controlifyTelegram"
)


type ControlTelegram struct {
	config *ControlTelegramConfig
	logger Logger
	Bot controlifyTelegram.Controller
}

type ControlTelegramConfig struct {
	Enabled bool `json:"enabled"`
	// Token - bot token.
	Token string `json:"token"`
	// AllowedUsers - bot accepts messages only from these users (id).
	AllowedUsers []int64 `json:"allowedUsers"`
	// AllowedChats - bot sends messages only to these chats (id).
	AllowedChats []int64 `json:"allowedChats"`
}

func (c *ControlTelegram) new(config *ControlTelegramConfig, logger Logger) {
	c.config = config
	c.logger = logger
	c.Bot = controlifyTelegram.Controller{}
	c.Bot.New(c)
}

func (c *ControlTelegram) SendFile(caption *string, filename string, reader io.Reader) {
	go c.Bot.SendFile(caption, filename, reader)
}

func (c *ControlTelegram) SendMessage(message string) {
	go c.Bot.SendMessage(message)
}

func (c *ControlTelegram) GetEnabled() bool {
	return c.config.Enabled
}

func (c *ControlTelegram) GetToken() string {
	return c.config.Token
}

func (c *ControlTelegram) LogError(message string) {
	go c.logger.Error(message)
}

func (c *ControlTelegram) GetAllowedChats() []int64 {
	return c.config.AllowedChats
}

func (c *ControlTelegram) GetAllowedUsers() []int64 {
	return c.config.AllowedUsers
}

