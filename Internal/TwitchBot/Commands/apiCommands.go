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

	"github.com/gempir/go-twitch-irc/v4"
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

func GetCurrentStreamInfo(broadcasterName string) (StreamInfo, error) {
	twApi := NewTwitchAPI()

	url := fmt.Sprintf("%s/streams?user_login=%s", twApi.BaseApiURL, broadcasterName)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return StreamInfo{}, err
	}

	req.Header.Set("Client-ID", twApi.ClientID)
	req.Header.Set("Authorization", "Bearer "+twApi.OAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return StreamInfo{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StreamInfo{}, err
	}

	var streamInfo StreamInfo
	err = json.Unmarshal(body, &streamInfo)
	if err != nil {
		return StreamInfo{}, err
	}

	if len(streamInfo.Data) == 0 {
		return StreamInfo{}, fmt.Errorf("streamer is offline or doesn't exist")
	}

	return streamInfo, nil
}

func GetCurrentGame(broadcasterName string) (string, error) {
	streamInfo, err := GetCurrentStreamInfo(broadcasterName)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("Current game is: %s", streamInfo.Data[0].GameName)

	return response, nil
}

func GetTitle(broadcasterName string) (string, error) {
	streamInfo, err := GetCurrentStreamInfo(broadcasterName)
	if err != nil {
		return "", err
	}

	//Maybe in the future, I"ll need some formatting
	response := streamInfo.Data[0].Title
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

func GetUserInfo(username string) (*twitch.User, error) {
	user, err := GetUserByLogin(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByLogin(username string) (*twitch.User, error) {
	twApi := NewTwitchAPI()
	url := fmt.Sprintf("%s/users?login=%s", twApi.BaseApiURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-ID", twApi.ClientID)
	req.Header.Set("Authorization", "Bearer "+twApi.OAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []struct {
			ID          string `json:"id"`
			Login       string `json:"login"`
			DisplayName string `json:"display_name"`
			Type        string `json:"type"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &twitch.User{
		ID:          response.Data[0].ID,
		Name:        response.Data[0].Login,
		DisplayName: response.Data[0].DisplayName,
	}, nil
}
