package main

import (
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	database "TelTwBot/Internal/Database"
	telegramBot "TelTwBot/Internal/Telegram"
	bot "TelTwBot/Internal/TwitchBot"
	"fmt"
	"log"
)

func main() {

	// Load database config
	dbConfigFile, err := config.ConfigPath(constants.DbConfigFile)
	if err != nil {
		log.Fatalf("Error getting database config path: %v", err)
	}

	dbConfig, err := config.LoadFromJSON[config.DbConfig](dbConfigFile)
	if err != nil {
		log.Fatalf("Error loading database config: %v", err)
	}

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	// Initialize database
	db, err := database.New(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	//Load telegramBot config
	tgBotFile, err := config.ConfigPath(constants.TgSettingsFile)
	if err != nil {
		log.Fatalf("Error getting Telegram config path: %v", err)
	}

	//Initialize telegramBot
	tgBot, err := telegramBot.NewTelegramNotifierFromConfigFile(tgBotFile)
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}

	rnd := config.InitRandom()

	//Load greetingFile
	greetFile, err := config.ConfigPath(constants.GreetingsFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	greeter, err := bot.NewGreeter(greetFile, rnd)

	if err != nil {
		log.Fatalf("Error while loading greetings file: %v", err)
	}

	log.Printf("%sLoaded %d greetings from file.", constants.Green, greeter.Count())

	//load duelFile
	duelFile, err := config.ConfigPath(constants.DuelsFile)
	if err != nil {
		log.Fatalf("Error gettings duels file path: %v", err)
	}

	allDuels, err := config.LoadFromJSON[[]config.DuelMsg](duelFile)
	if err != nil {
		log.Fatalf("Error while parsing duels file: %v", err)
	}

	//Initialize twitchBot
	twBot, err := bot.New(greeter, allDuels, tgBot)

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
