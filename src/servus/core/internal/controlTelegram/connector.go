package controlTelegram

import (
	"errors"
	"io"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// bridge between bot and controller.
type connector struct {
	bot        *bot
	controller *Controller
	commands   map[string]func(args []string)
}

// create new ControlTelegram instance.
func (c *connector) New(ctrl *Controller) error {
	// init.
	c.commands = make(map[string]func(args []string))
	c.controller = ctrl

	// bot.
	c.bot = &bot{}
	c.bot.OnUpdate(c.onUpdate)
	token := c.controller.GetToken()
	return c.bot.New(token)
}

// when message coming from user.
func (c *connector) onUpdate(update tgbotapi.Update) {
	c.fetchAllowedUsers(func(userID int64) bool {
		// disable if no commands / message.
		var isNoMessage = update.Message == nil || update.Message.From == nil
		if len(c.commands) < 1 || isNoMessage {
			return true
		}

		// check is message from allowed user.
		var isMessageFromAllowedUser = userID == update.Message.From.ID
		if !isMessageFromAllowedUser {
			return false
		}

		var message = update.Message.Text

		// split message by space, maybe it's command?
		var messageSpaced = strings.Split(message, " ")
		var command, ok = c.commands[messageSpaced[0]]

		// not a command.
		if !ok {
			return true
		}

		// execute callback (remove first element - its command).
		command(messageSpaced[1:])

		return true
	})
}

// add command. When user type command, executes callback.
//
// ex: /command arg1 arg2 arg3
func (c *connector) AddCommand(command string, callback func(args []string)) {
	c.commands[command] = callback
}

// send file to allowed chats.
func (c *connector) SendFile(caption *string, filename string, reader io.Reader) (err error) {
	if c.bot == nil || reader == nil {
		return errors.New("[telegram/sendfile]: nil bot or reader")
	}
	var file = File{}
	file.New(caption, filename, reader)
	c.fetchAllowedChats(func(chatID int64) bool {
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
	c.fetchAllowedChats(func(chatID int64) bool {
		_, err = c.bot.SendMessage(chatID, message)
		return err != nil
	})
	return
}

// executes callback on every allowed chat ID.
//
// Callback must return bool where true = stop.
func (c *connector) fetchAllowedChats(callback func(chatID int64) bool) {
	if callback == nil {
		return
	}
	var allowedChats = c.controller.GetAllowedChats()
	for index := range allowedChats {
		stop := callback(allowedChats[index])
		if stop {
			break
		}
	}
}

// executes callback on every allowed user.
//
// Callback must return bool where true = stop.
func (c *connector) fetchAllowedUsers(callback func(userID int64) bool) {
	if callback == nil {
		return
	}
	var allowedUsers = c.controller.GetAllowedUsers()
	for index := range allowedUsers {
		stop := callback(allowedUsers[index])
		if stop {
			break
		}
	}
}
