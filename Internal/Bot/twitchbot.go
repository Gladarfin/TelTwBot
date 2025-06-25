package bot

import (
	commands "TelTwBot/Internal/Commands"
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	telegramBot "TelTwBot/Internal/Telegram"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchBot struct {
	Client  *twitch.Client
	Greeter *Greeter
}

var tgBot *telegramBot.TelegramNotifier

func New(greeter *Greeter) (*TwitchBot, error) {
	tgBotFile, err := config.ConfigPath(constants.TgSettingsFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	tgBot, err = telegramBot.NewTelegramNotifierFromConfigFile(tgBotFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	tokenFilePath, err := config.ConfigPath(constants.TokenFile)
	if err != nil {
		log.Fatalf("Error getting token path: %v", err)
	}

	tokenData, err := os.ReadFile(tokenFilePath)
	if err != nil {
		log.Fatalf("Error reading token file: %v", err)
	}
	client := twitch.NewClient(constants.BotUsername, string(tokenData))

	return &TwitchBot{
		Client:  client,
		Greeter: greeter,
	}, nil
}

func (tb *TwitchBot) Connect() error {
	tb.Client.OnConnect(func() {
		log.Printf("%s✅ Bot connected to Twitch IRC!", constants.Blue)
		tgBot.SendMessage(fmt.Sprintf("[%s] ✅Bot connected to Twitch IRC!", time.Now().Format("15:04:05")))
		tb.Client.Join(constants.Channel)
	})
	tb.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {

		log.Printf("%s[%s] %s: %s\n", constants.White, message.Channel, message.User.Name, message.Message)

		//Return list of all commands for "!help"
		if strings.ToLower(message.Message) == "!help" {
			commandsList := commands.GetAllCommands()
			SayAndLog(tb.Client, constants.Channel, commandsList, constants.BotUsername)
		}
		//Generate random greeting for "!hello" command
		if strings.ToLower(message.Message) == "!hello" {
			greeting := tb.Greeter.GetRandomGreeting()
			response := fmt.Sprintf("@%s %s means 'hello' in %s", message.User.Name, greeting.Text, greeting.Language)
			SayAndLog(tb.Client, constants.Channel, response, constants.BotUsername)
		}
	})

	err := tb.Client.Connect()
	fmt.Println(err)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}

	return nil
}
