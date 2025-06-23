package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

const (
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"
)

//Greetings

type Greeting struct {
	Language string
	Text     string
}

var greetings []Greeting

// Preload on bot start
func loadGreetings(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			greetings = append(greetings, Greeting{
				Language: strings.TrimSpace(parts[0]),
				Text:     strings.TrimSpace(parts[1]),
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// Get random Greeting
func getRandomGreeting() Greeting {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return greetings[rand.Intn(len(greetings))]
}

func readTokenFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func sayAndLog(client *twitch.Client, channel string, message string, botUsername string) {
	log.Printf("%s [%s] answer to %s\n",
		Magenta,
		botUsername,
		message)
	client.Say(channel, message)
}

func main() {
	chanel := "gladarfin"
	botUsername := "gladarfin_bot"
	errGreetings := loadGreetings("hello.txt")

	if errGreetings != nil {
		log.Fatalf("Error while loading greetings: %v", errGreetings)
	}

	log.Printf("%sLoaded %d greetings from file.", Green, len(greetings))

	accessToken := readTokenFromFile(".client")
	client := twitch.NewClient("gladarfin_bot", accessToken)
	client.OnConnect(func() {
		log.Printf("%s✅ Bot connected to Twitch IRC!", Blue)
		client.Join(chanel)
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		//Log all messages from users to console
		log.Printf("%s[%s] %s: %s\n", White, message.Channel, message.User.Name, message.Message)

		//Responde to commands basic

		//Generate random greeting for "!hello" command
		if strings.ToLower(message.Message) == "!hello" {
			greeting := getRandomGreeting()
			response := fmt.Sprintf("@%s %s - 'hello' in %s", message.User.Name, greeting.Text, greeting.Language)
			sayAndLog(client, chanel, response, botUsername)
		}
	})

	err := client.Connect()
	fmt.Println(err)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}
}
