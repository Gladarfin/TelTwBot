package twBotCommands

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func GetUserRole(user *twitch.User) (string, error) {
	var roles []string

	if isBroadcaster(user) {
		roles = append(roles, "broadcaster")
	}
	if isModerator(user) {
		roles = append(roles, "moderator")
	}
	if isSubscriber(user) {
		roles = append(roles, "subscriber")
	}

	if len(roles) == 0 {
		return "", fmt.Errorf("user doesn't exist in current chatroom")
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("User %s is: ", user.Name))

	switch len(roles) {
	case 1:
		result.WriteString(roles[0] + ".")
	case 2:
		result.WriteString(fmt.Sprintf("%s and %s.", roles[0], roles[1]))
	default:
		allButLast := strings.Join(roles[:len(roles)-1], ", ")
		result.WriteString(fmt.Sprintf("%s, and %s.", allButLast, roles[len(roles)-1]))
	}

	return result.String(), nil
}

func isBroadcaster(user *twitch.User) bool {
	return hasBadge(user, "broadcaster")
}

func isModerator(user *twitch.User) bool {
	return hasBadge(user, "moderator") || isBroadcaster(user)
}

func isSubscriber(user *twitch.User) bool {
	return hasBadge(user, "subscriber")
}

func hasBadge(user *twitch.User, badge string) bool {
	if value, ok := user.Badges[badge]; ok {
		return value > 0
	}
	return false
}
