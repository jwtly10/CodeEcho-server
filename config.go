package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	OpenaiApiKey string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		OpenaiApiKey: os.Getenv("OPENAI_API_KEY"),
	}

	return config, nil
}
