package controlifyTelegram

import (
	"errors"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Controller - controls Telegram bot.
type Controller struct {
	outside Telegramer
	bot     *Bot
}

type Telegramer interface {
	// GetToken - get bot token.
	GetToken() string
	// GetAllowedChats - get chats ids where bot can send messages.
	GetAllowedChats() []int64
	// GetAllowedUsers - get users ids where bot can receive messages.
	GetAllowedUsers() []int64
}

// New - create new ControlTelegram instance.
func (c *Controller) New(t Telegramer) (err error) {
	c.outside = t
	c.bot = &Bot{}
	c.bot.OnUpdate(c.onUpdate)
	token := c.outside.GetToken()
	err = c.bot.New(token)
	return
}

// onUpdate - when message coming from user.
func (c *Controller) onUpdate(update tgbotapi.Update) {
	c.allowedUsersCallback(func(userID int64) bool {
		return userID == update.Message.From.ID
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
func (c *Controller) SendFile(caption *string, filename string, reader io.Reader) (err error) {
	if c.bot == nil || reader == nil {
		return errors.New("[telegram/sendfile]: nil bot or reader")
	}
	var file = File{}
	file.New(caption, filename, reader)
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err = c.bot.SendFile(chatID, &file)
		return err != nil
	})
	return
}

// SendMessage - send message to allowed chats.
func (c *Controller) SendMessage(message string) (err error) {
	if c.bot == nil {
		return
	}
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err = c.bot.SendMessage(chatID, message)
		return err != nil
	})
	return
}
