package config

import (
	constants "TelTwBot/Internal/Config/Constants"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type DuelMsg struct {
	AnnounceMessage string `json:"AnnounceMessage"`
	DuelMessage     string `json:"DuelMessage"`
	IsDraw          bool   `json:"IsDraw"`
}

func InitRandom() *rand.Rand {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	return rand.New(source)
}

func ReadTokenFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func ConfigPath(filename string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exeDir := filepath.Join(filepath.Dir(exePath), "../..")
	return filepath.Join(exeDir, constants.ConfigDir, filename), nil
}

func LoadDbConfig(path string) (*DbConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config DbConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadDuels(path string) ([]DuelMsg, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var duelMessages []DuelMsg
	if err := json.Unmarshal(file, &duelMessages); err != nil {
		return nil, err
	}
	return duelMessages, nil
}
