package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
	"github.com/joho/godotenv"
)

type HandlerEventFunc func(*ktwin.TwinEvent) error

func LoadEnv() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func StartServer(handleFuncTwin HandlerEventFunc) {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger()
		twinEvent := kevent.HandleRequest(r)
		if twinEvent == nil {
			logger.Error("Error handling cloud event request", nil)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error handling cloud event request"))
			return
		}

		logger.Info(fmt.Sprintf("Event TwinInstance: %s - Event TwinInterface: %s", twinEvent.TwinInstance, twinEvent.TwinInterface))

		if err := handleFuncTwin(twinEvent); err != nil {
			logger.Error("Error processing cloud event request", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error processing cloud event request"))
			return
		}
	}

	http.HandleFunc("/", handleFunc)

	logger := logger.NewLogger()
	logger.Info("Starting up server...")
	logger.Fatal("Server error", http.ListenAndServe(":8080", nil))
}
