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
	var config DbConfig
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
