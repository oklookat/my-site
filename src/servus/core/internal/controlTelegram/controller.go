package controlTelegram

import "io"

// writes information.
type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Panic(err error)
}

type Config struct {
	Enabled bool `json:"enabled"`
	// bot token.
	Token string `json:"token"`
	// bot accepts messages only from these users (id).
	AllowedUsers []int64 `json:"allowedUsers"`
	// bot sends messages only to these chats (id).
	AllowedChats []int64 `json:"allowedChats"`
}

type Controller struct {
	config *Config
	logger Logger
	bot    connector
}

func (c *Controller) New(config *Config, logger Logger) {
	if config == nil {
		panic("[control/telegram]: config nil pointer.")
	}
	if logger == nil {
		panic("[control/telegram]: logger nil pointer.")
	}
	c.config = config
	c.logger = logger
	c.bot = connector{}
	c.bot.New(c)
}

func (c *Controller) checkErr(err error) {
	if err != nil {
		c.logger.Error("[control/telegram]: " + err.Error())
	}
}

func (c *Controller) SendFile(caption *string, filename string, reader io.Reader) {
	err := c.bot.SendFile(caption, filename, reader)
	c.checkErr(err)
}

func (c *Controller) SendMessage(message string) {
	err := c.bot.SendMessage(message)
	c.checkErr(err)
}

func (c *Controller) GetToken() string {
	return c.config.Token
}

func (c *Controller) GetAllowedChats() []int64 {
	return c.config.AllowedChats
}

func (c *Controller) GetAllowedUsers() []int64 {
	return c.config.AllowedUsers
}
