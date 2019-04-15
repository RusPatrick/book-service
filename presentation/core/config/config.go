package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ruspatrick/book-service/domain/models"
)

const (
	ErrReadConfigText = "Failed read config"
)

var (
	config models.Config
)

func Read() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(ErrReadConfigText, err)
	}

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatal(ErrReadConfigText, err)
	}
}

func Get() models.Config {
	return config
}
