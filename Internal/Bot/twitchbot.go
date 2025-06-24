package bot

import (
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchBot struct {
	Client  *twitch.Client
	Greeter *Greeter
}

func New(greeter *Greeter) (*TwitchBot, error) {
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
		tb.Client.Join(constants.Chanel)
	})
	tb.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {

		log.Printf("%s[%s] %s: %s\n", constants.White, message.Channel, message.User.Name, message.Message)

		//Generate random greeting for "!hello" command
		if strings.ToLower(message.Message) == "!hello" {
			greeting := tb.Greeter.GetRandomGreeting()
			response := fmt.Sprintf("@%s %s means 'hello' in %s", message.User.Name, greeting.Text, greeting.Language)
			sayAndLog(tb.Client, constants.Chanel, response, constants.BotUsername)
		}
	})

	err := tb.Client.Connect()
	fmt.Println(err)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}

	return nil
}

func sayAndLog(client *twitch.Client, channel string, message string, botUsername string) {
	log.Printf("%s [%s] answer to %s\n",
		constants.Magenta,
		botUsername,
		message)
	client.Say(channel, message)
}
