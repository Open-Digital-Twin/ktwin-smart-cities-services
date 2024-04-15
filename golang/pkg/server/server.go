package server

import (
	"log"
	"net/http"
	"os"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func StartServer(handleFunc http.HandlerFunc) {
	logger := logger.NewLogger()

	http.HandleFunc("/", handleFunc)

	logger.Info("Starting up server...")
	logger.Fatal("Server error", http.ListenAndServe(":8080", nil))
}
