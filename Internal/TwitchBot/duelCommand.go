package bot

import (
	constants "TelTwBot/Internal/Config/Constants"
	"fmt"
	"time"
)

func (tb *TwitchBot) StartDuel(username string) {
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
		//Same user invoke duel twice
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
		SayAndLog(
			tb.Client,
			constants.Channel,
			fmt.Sprintf("⚔️DUEL! @%s has accepted @%s's challenge! Let the battle begin!", username, tb.CurrentDuel.Initiator),
			constants.BotUsername)

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
