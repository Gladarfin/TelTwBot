package config

import (
	constants "TelTwBot/Internal/Config/Constants"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

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

	exeDir := filepath.Join(filepath.Dir(exePath), "../../..")
	return filepath.Join(exeDir, constants.ConfigDir, filename), nil
}
