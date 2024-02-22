package main

import (
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/prerecorded"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	OpenaiApiKey   string
	DeepgramApiKey string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		OpenaiApiKey:   os.Getenv("OPENAI_API_KEY"),
		DeepgramApiKey: os.Getenv("DEEPGRAM_API_KEY"),
	}

	client.InitWithDefault()

	return config, nil
}
