package controlTelegram

import (
	"errors"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// bridge between lowlevel and controller.
type connector struct {
	outside telegramer
	bot     *bot
}

type telegramer interface {
	// get bot token.
	GetToken() string
	// get chats ids where bot can send messages.
	GetAllowedChats() []int64
	// get user IDs from which the bot can receive messages.
	GetAllowedUsers() []int64
}

// create new ControlTelegram instance.
func (c *connector) New(t telegramer) (err error) {
	c.outside = t
	c.bot = &bot{}
	c.bot.OnUpdate(c.onUpdate)
	token := c.outside.GetToken()
	err = c.bot.New(token)
	return
}

// when message coming from user.
func (c *connector) onUpdate(update tgbotapi.Update) {
	c.allowedUsersCallback(func(userID int64) bool {
		return userID == update.Message.From.ID
	})
}

// executes callback on every allowedChats ID. Callback must return bool where true = stop.
func (c *connector) allowedChatsCallback(callback func(chatID int64) bool) {
	if callback == nil {
		return
	}
	var allowedChats = c.outside.GetAllowedChats()
	for index := range allowedChats {
		stop := callback(allowedChats[index])
		if stop {
			break
		}
	}
}

// executes callback on every allowedChats ID. Callback must return bool where true = stop.
func (c *connector) allowedUsersCallback(callback func(userID int64) bool) {
	if callback == nil {
		return
	}
	var allowedUsers = c.outside.GetAllowedUsers()
	for index := range allowedUsers {
		stop := callback(allowedUsers[index])
		if stop {
			break
		}
	}
}

// send file to allowed chats.
func (c *connector) SendFile(caption *string, filename string, reader io.Reader) (err error) {
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

// send message to allowed chats.
func (c *connector) SendMessage(message string) (err error) {
	if c.bot == nil {
		return
	}
	c.allowedChatsCallback(func(chatID int64) bool {
		_, err = c.bot.SendMessage(chatID, message)
		return err != nil
	})
	return
}
