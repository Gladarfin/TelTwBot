package bot

import (
	constants "TelTwBot/Internal/Config/Constants"
	twBotCommands "TelTwBot/Internal/TwitchBot/Commands"
	"fmt"
	"log"
	"strconv"
	"strings"
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
				for _, msg := range commandsList {
					SayAndLog(tb.Client, constants.Channel, msg, constants.BotUsername)
				}
			},
		},
		{
			Name:        "!hello",
			Description: "Displays a random greeting to user.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				greeting := tb.Greeter.GetRandomGreeting()
				response := fmt.Sprintf("@%s, %s * means 'hello' in %s *", greeting.Text, message.User.Name, greeting.Language)
				SayAndLog(tb.Client, constants.Channel, response, constants.BotUsername)
				log.Printf("[%s] ✅Processed !hello command for %s.", time.Now().Format("15:04:05"), message.User.Name)
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
				log.Printf("[%s] ✅Processed !title command for %s.", time.Now().Format("15:04:05"), message.User.Name)
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
				log.Printf("[%s] ✅Processed !game command for %s.", time.Now().Format("15:04:05"), message.User.Name)
			},
		},
		{
			Name:        "!who",
			Description: "Shows participating streamers.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				friends, err := twBotCommands.GetStreamers()
				if err != nil {
					log.Printf("[%s]❌Failed to get streamers list.", time.Now().Format("15:04:05"))
				}
				SayAndLog(tb.Client, constants.Channel, friends, constants.BotUsername)
				log.Printf("[%s] ✅Processed !who command for %s.", time.Now().Format("15:04:05"), message.User.Name)
			},
		},
		{
			//There is a problem: that Helix API doesn't have info for user role on channel, so I need to think about how to work around it.
			//Now its command shows the roles for users that invoke it only.
			Name:        "!role",
			Description: "Shows the user role on current channel.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {

				args := strings.Fields(message.Message)
				var targetUser string
				if len(args) > 0 {
					//if username specified we use this username, if not we use username of user which invoke command
					targetUser = args[0]
				} else {
					targetUser = message.User.Name
				}

				roles, err := twBotCommands.GetUserRole(&message.User)
				if err != nil {
					log.Printf("[%s]❌%s has no special roles in this channel.", time.Now().Format("15:04:05"), message.User.Name)
					msg := fmt.Sprintf("%s has no special roles in this channel.", message.User.Name)
					SayAndLog(tb.Client, constants.Channel, msg, constants.BotUsername)
				}

				SayAndLog(tb.Client, constants.Channel, roles, constants.BotUsername)
				log.Printf("[%s] ✅Processed !role command for %s", time.Now().Format("15:04:05"), targetUser)
			},
		},
		{
			Name:        "!stats",
			Description: "Shows user stats.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				stats, err := twBotCommands.GetStats(message.User.Name)
				if err != nil {
					log.Printf("[%s]❌ Failed to get stats for %s: %v", time.Now().Format("15:04:05"), message.User.Name, err)
					SayAndLog(tb.Client, constants.Channel, "Sorry, couldn't retrieve your stats. Please try again later.", constants.BotUsername)
					return
				}
				SayAndLog(tb.Client, constants.Channel, stats, constants.BotUsername)
				log.Printf("[%s] ✅Processed !stats command for %s.", time.Now().Format("15:04:05"), message.User.Name)
			},
		},
		{
			Name:        "!duel",
			Description: "Starts the duel with other user.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				tb.StartDuel(message.User.Name)
			},
		},
		{
			Name:        "!up",
			Description: "Increase selected stat by 1 if there is enough free points.",
			Handler: func(tb *TwitchBot, message twitch.PrivateMessage) {
				args := strings.Fields(message.Message)

				if len(args) != 2 {
					log.Printf("[%s]❌Failed to increase stat for %s. The command contains an incorrect number of arguments.", time.Now().Format("15:04:05"), message.User.Name)
					SayAndLog(tb.Client, constants.Channel, "The stat command should contain the name of the stat that you want to increase and value: !up <stat_to_increse> <value>.", constants.BotUsername)
					return
				}

				val, err := strconv.Atoi(args[1])
				if err != nil {
					log.Printf("[%s]❌Error converting '%s' to int: %v", time.Now().Format("15:04:05"), args[1], err)
					SayAndLog(tb.Client, constants.Channel, "Second argument in command should be integer value.", constants.BotUsername)
					return
				}

				if val <= 0 {
					log.Printf("[%s]❌Error, value argument should be greater than 0.", time.Now().Format("15:04:05"))
					SayAndLog(tb.Client, constants.Channel, "Second argument in command should be greater than 0.", constants.BotUsername)
					return
				}

				if args[0] == "free-points" {
					log.Printf("[%s]❌Error, you can't increase free-points with this command.", time.Now().Format("15:04:05"))
					SayAndLog(tb.Client, constants.Channel, "Can't increase free-points stat this way.", constants.BotUsername)
					return
				}

				stats, err := twBotCommands.UpStat(message.User.Name, args[0], val)
				if err != nil {
					log.Printf("[%s]❌Failed to increase stat for %s: %s.", time.Now().Format("15:04:05"), message.User.Name, err)
					SayAndLog(tb.Client, constants.Channel, "Failed to increase stat.", constants.BotUsername)
					return
				}

				SayAndLog(tb.Client, constants.Channel, stats, constants.BotUsername)
				log.Printf("[%s] ✅Processed !up command for %s.", time.Now().Format("15:04:05"), message.User.Name)
			},
		},
	}
}
