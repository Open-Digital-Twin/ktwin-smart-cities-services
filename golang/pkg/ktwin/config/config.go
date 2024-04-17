package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if os.Getenv("ENV") == "test" {
		err := godotenv.Load("../local.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
