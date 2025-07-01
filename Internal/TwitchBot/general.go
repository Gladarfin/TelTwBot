package bot

import (
	constants "TelTwBot/Internal/Config/Constants"
	"log"

	"github.com/gempir/go-twitch-irc/v4"
)

func SayAndLog(client *twitch.Client, channel string, message string, botUsername string) {
	log.Printf("%s [%s] answer with '%s'\n",
		constants.Magenta,
		botUsername,
		message)
	client.Say(channel, message)
}
