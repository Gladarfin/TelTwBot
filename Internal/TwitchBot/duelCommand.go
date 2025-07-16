package bot

import (
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"
)

func (tb *TwitchBot) StartDuel(username string, duels []config.DuelMsg) {
	tb.DuelMutex.Lock()
	defer tb.DuelMutex.Unlock()

	curTime := time.Now()

	if tb.IsDuelCooldownActive && curTime.Sub(tb.LastDuelTime) < 5*time.Minute {
		remainingTime := time.Until(tb.LastDuelTime.Add(5 * time.Minute)).Round(time.Second)
		SayAndLog(
			tb.Client,
			constants.Channel,
			fmt.Sprintf("@%s, duels are on on cooldown. Please wait %s before challenging again.", username, remainingTime),
			constants.BotUsername)
		return
	}

	//If there is active duel
	if tb.CurrentDuel != nil && tb.CurrentDuel.IsActive {
		if tb.CurrentDuel.Initiator == username {
			SayAndLog(
				tb.Client,
				constants.Channel,
				fmt.Sprintf("%s, you've already challenged someone, wait for a response.", username),
				constants.BotUsername)
		}

		//Accept the duel
		tb.CurrentDuel.Challenger = username
		tb.CurrentDuel.Timer.Stop()
		tb.CurrentDuel.IsActive = false
		curDuel, winner, err := getDuel(duels)
		if err != nil {
			log.Fatalf("%s", err)
			SayAndLog(
				tb.Client,
				constants.Channel,
				"There is some error! Contact the administrator.",
				constants.BotUsername)
			return
		}
		formatedAnnounce := fmt.Sprintf(curDuel.AnnounceMessage, tb.CurrentDuel.Initiator, username)
		SayAndLog(
			tb.Client,
			constants.Channel,
			formatedAnnounce,
			constants.BotUsername)

		var formatedDuelMessage string
		switch winner {
		case 0:
			formatedDuelMessage = curDuel.DuelMessage

		case 1:
			formatedDuelMessage = fmt.Sprintf(curDuel.DuelMessage, tb.CurrentDuel.Initiator)
		case 2:
			formatedDuelMessage = fmt.Sprintf(curDuel.DuelMessage, username)
		}
		if winner >= 0 {
			SayAndLog(
				tb.Client,
				constants.Channel,
				formatedDuelMessage,
				constants.BotUsername)
		}
		//Set cooldown between duels
		tb.LastDuelTime = curTime
		tb.IsDuelCooldownActive = true
		tb.CurrentDuel = nil

		time.AfterFunc(5*time.Minute, func() {
			tb.DuelMutex.Lock()
			tb.IsDuelCooldownActive = false
			tb.DuelMutex.Unlock()
			SayAndLog(
				tb.Client,
				constants.Channel,
				"🔄The duel cooldown has ended! You can now challenge others again with !duel",
				constants.BotUsername)
		})
		return
	}

	//Clean expired duel
	if tb.CurrentDuel != nil && curTime.Sub(tb.CurrentDuel.CreationTime) >= time.Minute {
		tb.CurrentDuel = nil
	}

	tb.CurrentDuel = &DuelChallenge{
		Initiator:    username,
		IsActive:     true,
		CreationTime: curTime,
		Timer: time.AfterFunc(time.Minute, func() {
			tb.DuelMutex.Lock()
			defer tb.DuelMutex.Unlock()

			if tb.CurrentDuel != nil && tb.CurrentDuel.IsActive {
				SayAndLog(
					tb.Client,
					constants.Channel,
					fmt.Sprintf("@%s's duel challenge has expired with no takers.", tb.CurrentDuel.Initiator),
					constants.BotUsername)
			}
		}),
	}

	SayAndLog(
		tb.Client,
		constants.Channel,
		fmt.Sprintf("@%s has issued a duel challenge! Type !duel in the next 60 seconds to accept!", username),
		constants.BotUsername)
}

func getDuel(duels []config.DuelMsg) (config.DuelMsg, int, error) {
	newDuel, winner := getRandomDuel(duels)

	return newDuel, winner, nil
}

func getRandomDuel(duels []config.DuelMsg) (config.DuelMsg, int) {
	const MAX_DIFF = 5

	user1Result := rand.IntN(100)
	user2Result := rand.IntN(100)
	//Close duels == draw, so i set difference between rolls to 5 (MAX_DIFF)
	isCloseDuel := math.Abs(float64(user1Result)-float64(user2Result)) <= float64(MAX_DIFF)
	//res = 0 - draw (default), 1 - first win, 2 - second win
	res := 0
	if user1Result > user2Result {
		res = 1
	} else {
		res = 2
	}
	sortedDuels := getSortedDuels(duels, isCloseDuel)

	duelNumber := rand.IntN(len(sortedDuels))

	return sortedDuels[duelNumber], res
}

func getSortedDuels(allDuels []config.DuelMsg, isDraw bool) []config.DuelMsg {
	//i think for 30 records we don't need pre-allocation for array here, e.g.: make([]config.DuelMsg, 0, len(duels)/2)
	var duels []config.DuelMsg
	for _, duel := range allDuels {
		if duel.IsDraw == isDraw {
			duels = append(duels, duel)
		}
	}
	return duels
}
