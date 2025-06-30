package twBotCommands

import (
	config "TelTwBot/Internal/Config"
	constants "TelTwBot/Internal/Config/Constants"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type TwitchAPI struct {
	ClientID   string
	OAuthToken string
	BaseApiURL string
}

type StreamInfo struct {
	Data []struct {
		UserID      string `json:"user_id"`
		Username    string `json:"user_name"`
		GameID      string `json:"game_id"`
		GameName    string `json:"game_name"`
		Title       string `json:"title"`
		ViewerCount int    `json:"viewer_count"`
	} `json:"data"`
}

type ApiConfig struct {
	ClientID   string
	OAuthToken string
}

func NewTwitchAPI() *TwitchAPI {
	helixFile, err := config.ConfigPath(constants.HelixFile)
	if err != nil {
		log.Fatalf("Error getting Helix config path: %v", err)
	}

	conf, err := LoadConfigFromFile(helixFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return nil
	}

	return &TwitchAPI{
		ClientID:   conf.ClientID,
		OAuthToken: conf.OAuthToken,
		BaseApiURL: "https://api.twitch.tv/helix",
	}
}

func GetCurrentGame(broadcasterName string) (string, error) {
	twApi := NewTwitchAPI()

	url := fmt.Sprintf("%s/streams?user_login=%s", twApi.BaseApiURL, broadcasterName)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Client-ID", twApi.ClientID)
	req.Header.Set("Authorization", "Bearer "+twApi.OAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var streamInfo StreamInfo
	err = json.Unmarshal(body, &streamInfo)
	if err != nil {
		return "", err
	}

	if len(streamInfo.Data) == 0 {
		return "", fmt.Errorf("Streamer is offline or doesn't exist.")
	}

	response := fmt.Sprintf("Current game is: %s", streamInfo.Data[0].GameName)

	return response, nil
}

func LoadConfigFromFile(filename string) (*ApiConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ApiConfig{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "clientID":
			config.ClientID = value
		case "oauth":
			config.OAuthToken = value
		}
	}

	if config.ClientID == "" {
		return nil, fmt.Errorf("clientID not found in config file")
	}
	if config.OAuthToken == "" {
		return nil, fmt.Errorf("oauth not found in config file")
	}

	return config, nil
}
