package twBotCommands

import (
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	"bufio"
	"log"
	"os"
	"strings"
)

type Friend struct {
	Nickname string
	Channel  string
}

func GetStreamers() (string, error) {
	frFile, err := config.ConfigPath(constants.FriendsFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	file, err := os.Open(frFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var friendsText strings.Builder
	friendsText.WriteString("Our friends (for tonight): ")
	friendsText.WriteString(" ---------------------------------------------- ")
	for scanner.Scan() {
		friendsText.WriteString(scanner.Text())
		friendsText.WriteString("       ")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return friendsText.String(), nil
}
