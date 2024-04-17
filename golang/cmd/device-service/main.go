package main

import (
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/device-service/service"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/config"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/server"
)

func main() {
	config.LoadEnv()
	server.StartServer(service.HandleEvent)
}
