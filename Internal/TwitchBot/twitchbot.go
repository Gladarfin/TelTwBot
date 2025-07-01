package bot

import (
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	botInterfaces "TelTwBot/Internal/Interfaces"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchBot struct {
	Client     *twitch.Client
	Greeter    *Greeter
	startTime  time.Time
	streamLive bool
	tgBot      botInterfaces.TelegramNotifierInterface
	commands   []Command
}

type Command struct {
	Name        string
	Description string
	Handler     func(tb *TwitchBot, message twitch.PrivateMessage)
}

var _ botInterfaces.TwitchBotInterface = (*TwitchBot)(nil)

func New(greeter *Greeter, tgNotifier botInterfaces.TelegramNotifierInterface) (*TwitchBot, error) {
	tokenFilePath, err := config.ConfigPath(constants.TokenFile)
	if err != nil {
		log.Fatalf("Error getting token path: %v", err)
	}

	tokenData, err := os.ReadFile(tokenFilePath)
	if err != nil {
		log.Fatalf("Error reading token file: %v", err)
	}
	client := twitch.NewClient(constants.BotUsername, string(tokenData))
	client.SetIRCToken(string(tokenData))

	return &TwitchBot{
		Client:     client,
		Greeter:    greeter,
		startTime:  time.Now(),
		streamLive: false,
		tgBot:      tgNotifier,
	}, nil
}

func (tb *TwitchBot) Connect() error {
	tb.InitCommands()
	tb.Client.OnConnect(func() {
		log.Printf("%s✅Bot connected to Twitch IRC!", constants.Blue)
		tb.tgBot.SendMessage(fmt.Sprintf("[%s] ✅Bot connected to Twitch IRC!", time.Now().Format("15:04:05")))
		tb.Client.Join(constants.Channel)
		tb.streamLive = true
		tb.startTime = time.Now()
	})
	tb.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {

		cmdInput := strings.ToLower(message.Message)
		for _, cmd := range tb.commands {
			if cmdInput == cmd.Name {
				cmd.Handler(tb, message)
				break
			}
		}
		log.Printf("%s[%s] %s: %s\n", constants.White, message.Channel, message.User.Name, message.Message)
	})

	tb.Client.OnUserPartMessage(func(message twitch.UserPartMessage) {
		if message.User == constants.Channel {
			tb.streamLive = false
			log.Printf("Stream went offline at %s", time.Now().Format("15:04:05"))
			log.Printf("Trying to reconnect...")
			ReconnectTwitch(tb, 10)
		}
	})

	err := tb.Client.Connect()
	fmt.Println(err)
	if err != nil {
		log.Fatal("❌Failed to connect:", err)
		tb.tgBot.SendMessage(fmt.Sprintf("[%s] ❌Failed to connect: %s", time.Now().Format("15:04:05"), err))
	}

	return nil
}

func ReconnectTwitch(tb *TwitchBot, maxRetries int) {
	retryCount := 0
	baseDelay := 5 * time.Second

	for {
		err := tb.Client.Connect()
		if err == nil {
			retryCount = 0
			baseDelay = 5 * time.Second
			continue
		}

		log.Printf("❌Connection error: %v (retry %d/%d)", err, retryCount+1, maxRetries)
		message := fmt.Sprintf("❌Connection error: %v (retry %d/%d)", err, retryCount+1, maxRetries)
		tb.tgBot.SendMessage(message)

		if retryCount >= maxRetries {
			log.Fatal("❌Connection failed! Max retries reached.")
			tb.tgBot.SendMessage("❌Connection failed! Max retries reached.")
		}

		delay := baseDelay * time.Duration(1<<retryCount)
		delay += time.Duration(rand.Intn(2000)) * time.Millisecond

		log.Printf("Reconnecting in %v...", delay)
		time.Sleep(delay)
		retryCount++
	}
}

func (tb *TwitchBot) GetStreamUptime() (string, error) {
	if !tb.streamLive {
		return "🔴Stream is currently offline.", nil
	}

	uptime := time.Since(tb.startTime)

	return fmt.Sprintf("%02d:%02d:%02d",
		int(uptime.Hours()),
		int(uptime.Minutes())%60,
		int(uptime.Seconds())%60), nil
}

func GetAllCommands(tb *TwitchBot) string {
	var helpText strings.Builder
	helpText.WriteString("Available commands: ")
	helpText.WriteString(" --------------------------------------- ")
	//Yes, it's the year 2025 A.D., and we don't have multiline messages in Twitch.
	for _, cmd := range tb.commands {
		helpText.WriteString(fmt.Sprintf("%s: %s", cmd.Name, cmd.Description))
		helpText.WriteString(" --------------------------------------- ")
	}
	return helpText.String()
}
