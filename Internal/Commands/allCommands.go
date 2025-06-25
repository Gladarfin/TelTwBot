package commands

import (
	"fmt"
	"strings"
)

type Command struct {
	Name        string
	Description string
}

var allCommands = []Command{
	{
		Name:        "!help",
		Description: "Displays a list of available commands.",
	},
	{
		Name:        "!hello",
		Description: "Displays a random greeting to user.",
	},
}

func GetAllCommands() string {
	var helpText strings.Builder
	helpText.WriteString("Available commands: ")
	//Yes, it's the year 2025 A.D., and we don't have multiline messages in Twitch.
	helpText.WriteString(" ---------------------------------------------- ")
	for _, cmd := range allCommands {
		helpText.WriteString(fmt.Sprintf("%s: %s", cmd.Name, cmd.Description))
		helpText.WriteString(" ---------------------------------------------- ")
	}
	return helpText.String()
}
