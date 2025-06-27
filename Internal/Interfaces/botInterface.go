package interfaces

type TwitchBotInterface interface {
	GetStreamUptime() (string, error)
}

type TelegramNotifierInterface interface {
	SendMessage(text string) error
}
