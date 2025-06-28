package telegramBot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetBotCommands() []tgbotapi.BotCommand {
	return []tgbotapi.BotCommand{
		{Command: "uptime", Description: "Get stream uptime"},
		{Command: "haha", Description: "Just for test"},
		{Command: "help", Description: "Show help"},
	}
}
