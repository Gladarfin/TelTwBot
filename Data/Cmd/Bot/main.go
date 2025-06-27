package main

import (
	bot "TelTwBot/Internal/Bot"
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	telegramBot "TelTwBot/Internal/Telegram"
	"log"
)

func main() {
	tgBotFile, err := config.ConfigPath(constants.TgSettingsFile)
	if err != nil {
		log.Fatalf("Error getting Telegram config path: %v", err)
	}

	tgBot, err := telegramBot.NewTelegramNotifierFromConfigFile(tgBotFile)
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}

	rnd := config.InitRandom()
	greetFile, err := config.ConfigPath(constants.GreetingsFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	greeter, err := bot.NewGreeter(greetFile, rnd)

	if err != nil {
		log.Fatalf("Error while loading greetings file: %v", err)
	}

	log.Printf("%sLoaded %d greetings from file.", constants.Green, greeter.Count())

	twBot, err := bot.New(greeter, tgBot)

	if err != nil {
		log.Fatalf("Error creating bot %v", err)
	}

	go func() {
		if err := twBot.Connect(); err != nil {
			log.Fatalf("Twitch bot connection error: %v", err)
		}
	}()

	tgBot.StartListening(twBot)

	select {}
}
