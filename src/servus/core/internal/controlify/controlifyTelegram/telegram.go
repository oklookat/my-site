package controlify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
)

// ControlTelegram - controls Telegram bot.
type ControlTelegram struct {
	outside controlTelegrammer
	api  *Telegram
}

type controlTelegrammer interface {
	// getEnabled - is bot enabled.
	getEnabled() bool
	// getToken - get bot token.
	getToken() string
	// logError - error logging.
	logError(message string)
	// getAllowedChats - get chats ids where bot can send messages.
	getAllowedChats() []int64
	// getAllowedChats - get users ids where bot can receive messages.
	getAllowedUsers() []int64
}

// New - create new ControlTelegram instance.
func (c *ControlTelegram) New(t controlTelegrammer) {
	c.outside = t
	var isEnabled = t.getEnabled()
	if isEnabled {
		c.api = &Telegram{}
		c.api.OnUpdate(c.onUpdate)
		token := c.outside.getToken()
		err := c.api.New(token)
		if err != nil {
			c.outside.logError("[control/telegram] error while booting: " + err.Error())
		}
	}
}

// onUpdate - when message coming from user.
func (c *ControlTelegram) onUpdate(update tgbotapi.Update) {
	c.allowedUsersCallback(func(userID int64) bool {
		if userID != update.Message.From.ID {
			// not allowed user
			return false
		}
		// allowed user
		return true
	})
}

// sendMessage - send file to allowed chats.
func (c *ControlTelegram) sendFile(caption *string, filename string, reader io.Reader) {
	if c.api == nil {
		return
	}
	var file = TelegramFile{}
	file.New(caption, filename, reader)
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err := c.api.SendFile(chatID, &file)
		if err != nil {
			c.outside.logError("[control/telegram] error while sending file: " + err.Error())
			return true
		}
		return false
	})
}

// sendMessage - send message to allowed chats.
func (c *ControlTelegram) sendMessage(message string) {
	if c.api == nil {
		return
	}
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err := c.api.SendMessage(chatID, message)
		if err != nil {
			c.outside.logError("[control/telegram] error while sending message: " + err.Error())
			return true
		}
		return false
	})
}

// allowedChatsCallback - executes callback on every allowedChats ID. Callback must return bool where true = stop.
func (c *ControlTelegram) allowedChatsCallback(callback func(chatID int64) bool) {
	var allowedChats = c.outside.getAllowedChats()
	for index := range allowedChats {
		stop := callback(allowedChats[index])
		if stop {
			break
		}
	}
}

// allowedChatsCallback - executes callback on every allowedChats ID. Callback must return bool where true = stop.
func (c *ControlTelegram) allowedUsersCallback(callback func(userID int64) bool) {
	var allowedUsers = c.outside.getAllowedUsers()
	for index := range allowedUsers {
		stop := callback(allowedUsers[index])
		if stop {
			break
		}
	}
}
