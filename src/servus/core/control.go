package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"servus/core/internal/controlify"
)

type Control struct {
	core *Core
	tg controlTelegram
}

type controlTelegram struct {
	core *Core
	api  controlify.Telegram
}

func (n *Control) boot(c *Core) {
	n.core = c
	var isTG = n.core.Config.Control.Telegram.Enabled
	if isTG {
		var token = n.core.Config.Control.Telegram.Token
		n.tg = controlTelegram{}
		n.tg.core = c
		n.tg.api = controlify.Telegram{}
		n.tg.api.OnUpdate(n.tg.onUpdate)
		err := n.tg.api.New(token)
		if err != nil {
			n.core.Logger.Error("notify/telegram error: " + err.Error())
		}
	}
}

func (n *controlTelegram) onUpdate(update tgbotapi.Update) {
	var allowedUsers = n.core.Config.Control.Telegram.AllowedUsers
	for index := range allowedUsers {
		if allowedUsers[index] == update.Message.From.ID {
			break
		}
		return
	}
}

func (n *controlTelegram) sendFile(caption  *string, filename string, reader io.Reader) {
	var allowedChats = n.core.Config.Control.Telegram.AllowedChats
	var file = controlify.TelegramFile{}
	file.New(caption, filename, reader)
	for index := range allowedChats {
		_, _ = n.api.SendFile(allowedChats[index], &file)
	}
}
