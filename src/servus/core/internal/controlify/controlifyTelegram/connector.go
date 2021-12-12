package controlifyTelegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
)

// Bot - bot.
type Bot struct {
	bot      *tgbotapi.BotAPI
	onUpdate func(update tgbotapi.Update)
}

// New - create new instance.
func (t *Bot) New(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	t.bot = bot
	go t.watchUpdates()
	return err
}

// OnUpdate - when message coming from user.
func (t *Bot) OnUpdate(callback func(update tgbotapi.Update)) {
	t.onUpdate = callback
}

// watchUpdates - watch chat changes.
func (t *Bot) watchUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if t.onUpdate != nil {
			t.onUpdate(update)
		}
	}
}

func (t *Bot) SendMessage(chatID int64, message string) (tgbotapi.Message, error) {
	var req = tgbotapi.NewMessage(chatID, message)
	return t.bot.Send(req)
}

func (t *Bot) SendFile(chatID int64, file *File) (tgbotapi.Message, error) {
	var req = tgbotapi.NewDocument(chatID, file)
	if file.caption != nil {
		req.Caption = *file.caption
	}
	return t.bot.Send(req)
}

type File struct {
	caption  *string
	filename string
	reader   io.Reader
}

func (t *File) New(caption *string, filename string, reader io.Reader) {
	t.caption = caption
	t.filename = filename
	t.reader = reader
}

func (t *File) NeedsUpload() bool {
	return true
}

func (t *File) UploadData() (string, io.Reader, error) {
	return t.filename, t.reader, nil
}

func (t *File) SendData() string {
	return t.filename
}
