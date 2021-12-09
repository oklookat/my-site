package controlify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
)

type Telegram struct {
	bot      *tgbotapi.BotAPI
	onUpdate func(update tgbotapi.Update)
}

type TelegramContent struct {
	Message *string
	File    *TelegramFile
}

func (t *Telegram) New(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	t.bot = bot
	go t.watchUpdates()
	return err
}

func (t *Telegram) OnUpdate(callback func(update tgbotapi.Update)) {
	t.onUpdate = callback
}

func (t *Telegram) watchUpdates() {
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

func (t *Telegram) SendFile(chatID int64, file *TelegramFile) (tgbotapi.Message, error) {
	var req = tgbotapi.NewDocument(chatID, file)
	if file.caption != nil {
		req.Caption = *file.caption
	}
	return t.bot.Send(req)
}

type TelegramFile struct {
	caption  *string
	filename string
	reader   io.Reader
}

func (t *TelegramFile) New(caption *string, filename string, reader io.Reader) {
	t.caption = caption
	t.filename = filename
	t.reader = reader
}

func (t *TelegramFile) NeedsUpload() bool {
	return true
}

func (t *TelegramFile) UploadData() (string, io.Reader, error) {
	return t.filename, t.reader, nil
}

func (t *TelegramFile) SendData() string {
	return t.filename
}
