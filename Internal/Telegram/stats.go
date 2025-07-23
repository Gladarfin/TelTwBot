package telegramBot

import (
	db "TelTwBot/Internal/Database"
	"context"
	"fmt"
	"strings"
)

func GetStats(username string) (string, error) {
	database := db.GetInstance()
	stats, err := database.GetTwitchUserStats(context.Background(), username)
	if err != nil {
		if isUserNotFoundError(err) {
			return fmt.Sprintf("User %s not found in the database.", username), nil
		}
		return "", err
	}

	emojiMap := map[string]string{
		"strength":     "💪",
		"perception":   "👀",
		"endurance":    "🛡️",
		"charisma":     "🎭",
		"intelligence": "🧠",
		"agility":      "🏃‍♂️",
		"luck":         "🍀",
		"free-points":  "✨",
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("%s's stats: \n", username))
	for _, stat := range stats {
		emoji := emojiMap[stat.StatType]
		message.WriteString(fmt.Sprintf("%s %s : %d \n", emoji, stat.StatType, stat.Value))
	}
	return message.String(), nil
}

func isUserNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "user not found") ||
		strings.Contains(err.Error(), "no rows in result set")
}
