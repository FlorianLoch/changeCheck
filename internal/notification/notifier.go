package notification

import (
	"fmt"
	u "net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Nofifier interface {
	Notify(url *u.URL, debounce bool) error
}

type TelegramNotifier struct {
	bot       *tgbotapi.BotAPI
	chatID    int64
	debouncer Debouncer
}

func NewTelegramNotifier(botToken string, chatID int64, debouncer Debouncer) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{
		bot:       bot,
		chatID:    chatID,
		debouncer: debouncer,
	}, nil
}

func (t *TelegramNotifier) Notify(url *u.URL, shallDebounce bool) error {
	shallNotify, relayURL := t.debouncer.ShallNotify(url)
	if shallDebounce && shallNotify {
		return nil
	}

	msg := tgbotapi.NewMessage(t.chatID, fmt.Sprintf("'%s' has changed: %s", url, relayURL))

	_, err := t.bot.Send(msg)

	return err
}
