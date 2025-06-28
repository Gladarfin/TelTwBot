package telegramBot

import (
	botInterfaces "TelTwBot/Internal/Interfaces"
	helpers "TelTwBot/Internal/Telegram/Helpers"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

type BotConfig struct {
	BotToken string
	ChatID   int64
}

func NewTelegramNotifier(botToken string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	setBotCommands(bot)

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func NewTelegramNotifierFromConfigFile(filename string) (*TelegramNotifier, error) {
	config, err := LoadConfigFromFile(filename)
	if err != nil {
		return nil, err
	}
	return NewTelegramNotifier(config.BotToken, config.ChatID)
}

func setBotCommands(bot *tgbotapi.BotAPI) error {
	commands := GetBotCommands()
	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(config)
	return err
}

func LoadConfigFromFile(filename string) (*BotConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &BotConfig{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "botToken":
			config.BotToken = value
		case "chatId":
			id, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid chatId: %v", err)
			}
			config.ChatID = id
		}
	}

	if config.BotToken == "" {
		return nil, fmt.Errorf("botToken not found in config file")
	}
	if config.ChatID == 0 {
		return nil, fmt.Errorf("chatId not found in config file")
	}

	return config, nil
}

func (tn *TelegramNotifier) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(tn.chatID, text)
	_, err := tn.bot.Send(msg)
	return err
}

func (tn *TelegramNotifier) StartListening(twitchBot botInterfaces.TwitchBotInterface) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	updates := tn.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			tn.handleCommand(update, twitchBot)
		} else {
			tn.handleMessage(update)
		}
	}
}

func (tn *TelegramNotifier) handleCommand(update tgbotapi.Update, twitchBot botInterfaces.TwitchBotInterface) {
	command := strings.ToLower(update.Message.Command())
	args := update.Message.CommandArguments()

	switch command {
	case "uptime":
		tn.handleUptimeCommand(update, twitchBot)
	case "help":
		tn.handleHelpCommand(update)
	case "math":
		tn.handleMathCommand(update, args)
	default:
		tn.sendMessage(update.Message.Chat.ID, "Unknown command. Try /help")
	}
}

func (tn *TelegramNotifier) handleUptimeCommand(update tgbotapi.Update, twitchBot botInterfaces.TwitchBotInterface) {
	uptime, err := twitchBot.GetStreamUptime()
	if err != nil {
		tn.sendMessage(update.Message.Chat.ID, "Error checking uptime: "+err.Error())
		return
	}

	response := fmt.Sprintf("🕒 Stream uptime: %s", uptime)
	tn.sendMessage(update.Message.Chat.ID, response)
}

func (tn *TelegramNotifier) handleHelpCommand(update tgbotapi.Update) {
	commands := GetBotCommands()
	var helpText strings.Builder
	helpText.WriteString("🤖 <b>Available Commands:</b>\n\n")

	for _, cmd := range commands {
		helpText.WriteString(fmt.Sprintf("/%s - %s \n\n", cmd.Command, cmd.Description))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText.String())
	msg.ParseMode = "HTML"
	tn.bot.Send(msg)
}

func (tn *TelegramNotifier) handleMessage(update tgbotapi.Update) {

}

func (tn *TelegramNotifier) handleMathCommand(update tgbotapi.Update, args string) {
	if args == "" {
		tn.sendMessage(update.Message.Chat.ID, "Incorrect input. Usage: /math <expression> (e.g. /math 2+2)")
		return
	}

	result, err := helpers.EvaluateExpression(args)
	if err != nil {
		tn.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Error: %s\nUsage: /math <num> <op> <num> (e.g., 2 + 2)", err.Error()))
		return
	}

	response := fmt.Sprintf("🧮Result: %s = %v", args, result)
	tn.sendMessage(update.Message.Chat.ID, response)
}

func (tn *TelegramNotifier) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := tn.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending Telegram message: %v", err)
	}
}
