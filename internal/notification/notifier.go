package notification

type Nofifier interface {
	Notify(url string)
}

type TelegramNotifier struct {
}

func (t *TelegramNotifier) Notify(url string) {
	// TODO: implement
}
