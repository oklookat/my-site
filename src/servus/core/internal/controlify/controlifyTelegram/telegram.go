package controlifyTelegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
)

// Controller - controls Telegram bot.
type Controller struct {
	outside Telegramer
	bot     *Bot
}

type Telegramer interface {
	// GetEnabled - is bot enabled.
	GetEnabled() bool
	// GetToken - get bot token.
	GetToken() string
	// LogError - error logging.
	LogError(message string)
	// GetAllowedChats - get chats ids where bot can send messages.
	GetAllowedChats() []int64
	// GetAllowedUsers - get users ids where bot can receive messages.
	GetAllowedUsers() []int64
}

// New - create new ControlTelegram instance.
func (c *Controller) New(t Telegramer) {
	c.outside = t
	var isEnabled = t.GetEnabled()
	if isEnabled {
		c.bot = &Bot{}
		c.bot.OnUpdate(c.onUpdate)
		token := c.outside.GetToken()
		err := c.bot.New(token)
		if err != nil {
			c.outside.LogError("[control/telegram] error while booting: " + err.Error())
		}
	}
}

// onUpdate - when message coming from user.
func (c *Controller) onUpdate(update tgbotapi.Update) {
	c.allowedUsersCallback(func(userID int64) bool {
		if userID != update.Message.From.ID {
			// not allowed user
			return false
		}
		// allowed user
		return true
	})
}

// allowedChatsCallback - executes callback on every allowedChats ID. Callback must return bool where true = stop.
func (c *Controller) allowedChatsCallback(callback func(chatID int64) bool) {
	var allowedChats = c.outside.GetAllowedChats()
	for index := range allowedChats {
		stop := callback(allowedChats[index])
		if stop {
			break
		}
	}
}

// allowedChatsCallback - executes callback on every allowedChats ID. Callback must return bool where true = stop.
func (c *Controller) allowedUsersCallback(callback func(userID int64) bool) {
	var allowedUsers = c.outside.GetAllowedUsers()
	for index := range allowedUsers {
		stop := callback(allowedUsers[index])
		if stop {
			break
		}
	}
}

// SendFile - send file to allowed chats.
func (c *Controller) SendFile(caption *string, filename string, reader io.Reader) {
	if c.bot == nil {
		return
	}
	var file = File{}
	file.New(caption, filename, reader)
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err := c.bot.SendFile(chatID, &file)
		if err != nil {
			c.outside.LogError("[control/telegram] error while sending file: " + err.Error())
			return true
		}
		return false
	})
}

// SendMessage - send message to allowed chats.
func (c *Controller) SendMessage(message string) {
	if c.bot == nil {
		return
	}
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err := c.bot.SendMessage(chatID, message)
		if err != nil {
			c.outside.LogError("[control/telegram] error while sending message: " + err.Error())
			return true
		}
		return false
	})
}
