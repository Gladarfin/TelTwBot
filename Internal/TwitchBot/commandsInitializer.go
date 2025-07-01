package bot

import (
	constants "TelTwBot/Internal/Config/Constants"
	twBotCommands "TelTwBot/Internal/TwitchBot/Commands"
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

func (tb *TwitchBot) InitCommands() {
	tb.commands = []Command{
		{
			Name:        "!help",
			Description: "Displays a list of available commands.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				commandsList := GetAllCommands(tb)
				SayAndLog(tb.Client, constants.Channel, commandsList, constants.BotUsername)
			},
		},
		{
			Name:        "!hello",
			Description: "Displays a random greeting to user.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				greeting := tb.Greeter.GetRandomGreeting()
				response := fmt.Sprintf("@%s %s means 'hello' in %s", message.User.Name, greeting.Text, greeting.Language)
				SayAndLog(tb.Client, constants.Channel, response, constants.BotUsername)
			},
		},
		{
			Name:        "!title",
			Description: "Displays the current stream title.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				title, err := twBotCommands.GetTitle(constants.Channel)
				if err != nil {
					log.Printf("[%s]❌Failed to get title of the stream. Error: %s", time.Now().Format("15:04:05"), err)
				}
				SayAndLog(tb.Client, constants.Channel, title, constants.BotUsername)
			},
		},
		{
			Name:        "!game",
			Description: "Shows what game is currently being played.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				game, err := twBotCommands.GetCurrentGame(constants.Channel)
				if err != nil {
					log.Printf("[%s]❌Failed to get game name. Error: %s", time.Now().Format("15:04:05"), err)
				}
				SayAndLog(tb.Client, constants.Channel, game, constants.BotUsername)
			},
		},
		{
			Name:        "!who",
			Description: "Shows participating streamers",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				friends, err := twBotCommands.GetStreamers()
				if err != nil {
					log.Printf("[%s]❌Failed to get streamers list.", time.Now().Format("15:04:05"))
				}
				SayAndLog(tb.Client, constants.Channel, friends, constants.BotUsername)
			},
		}}
}
