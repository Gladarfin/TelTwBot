package twBotCommands

import (
	db "TelTwBot/Internal/Database"
	"context"
	"fmt"
	"strings"
)

func GetStats(username string) (string, error) {
	database := db.GetInstance()
	stats, err := database.GetOrCreateUserStats(context.Background(), username)
	if err != nil {
		return "", err
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("%s's stats: ", username))
	for _, stat := range stats {
		message.WriteString(fmt.Sprintf("%s: %d | ", stat.StatType, stat.Value))
	}
	return message.String(), nil
}

func UpStat(username string, stat string, val int) (string, error) {
	database := db.GetInstance()
	statName, newStatValue, err := database.UpdateUserStat(context.Background(), username, stat, val)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s's %s is now %d", username, statName, newStatValue), nil
}
