package notification

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Nofifier interface {
	Notify(url string) error
}

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotifier(botToken string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t *TelegramNotifier) Notify(url string) error {
	msg := tgbotapi.NewMessage(t.chatID, fmt.Sprintf("'%s' has changed.", url))

	_, err := t.bot.Send(msg)

	return err
}
