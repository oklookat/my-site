package controlTelegram

import (
	"errors"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// bot - bot.
type bot struct {
	api      *tgbotapi.BotAPI
	onUpdate func(update tgbotapi.Update)
}

// New - create new instance.
func (t *bot) New(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	t.api = bot
	go t.watchUpdates()
	return err
}

// OnUpdate - when message coming from user.
func (t *bot) OnUpdate(callback func(update tgbotapi.Update)) {
	t.onUpdate = callback
}

// watchUpdates - watch chat changes.
func (t *bot) watchUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if t.onUpdate != nil {
			t.onUpdate(update)
		}
	}
}

func (t *bot) SendMessage(chatID int64, message string) (tgbotapi.Message, error) {
	var req = tgbotapi.NewMessage(chatID, message)
	return t.api.Send(req)
}

func (t *bot) SendFile(chatID int64, file *File) (tgbotapi.Message, error) {
	if file == nil {
		return tgbotapi.Message{}, errors.New("[control/telegram]: empty file")
	}
	var req = tgbotapi.NewDocument(chatID, file)
	if file.caption != nil {
		req.Caption = *file.caption
	}
	return t.api.Send(req)
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
