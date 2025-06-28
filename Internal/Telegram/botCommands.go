package telegramBot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetBotCommands() []tgbotapi.BotCommand {
	return []tgbotapi.BotCommand{
		{Command: "uptime", Description: "Get stream uptime"},
		{Command: "test", Description: "Just for test"},
		{Command: "math", Description: "Do simple math (e.g. a + b)"},
		{Command: "help", Description: "Show help"},
	}
}
