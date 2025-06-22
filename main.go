package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gempir/go-twitch-irc/v4"
)

func readTokenFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func main() {
	accessToken := readTokenFromFile(".client")
	client := twitch.NewClient("gladarfin_bot", accessToken)

	client.OnConnect(func() {
		fmt.Println("✅ Bot connected to Twitch IRC!")
		client.Join("gempir")
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.Name, message.Message)
	})

	err := client.Connect()
	fmt.Println(err)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}

}
