package main

import (
	bot "TelTwBot/Internal/Bot"
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	"log"
)

func main() {

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

	twBot, err := bot.New(greeter)

	if err != nil {
		log.Fatalf("Error creating bot %v", err)
	}

	if err := twBot.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
}
