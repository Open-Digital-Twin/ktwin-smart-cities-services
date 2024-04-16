package main

import (
	"fmt"
	"net/http"

	parkingModel "github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-spot-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/server"
)

var logger = log.NewLogger()
var twinGraph ktwin.TwinGraph

func requestHandler(w http.ResponseWriter, r *http.Request) {
	twinEvent := ktwin.HandleRequest(r)

	if twinEvent == nil {
		logger.Error("Error handling cloud event request", nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error handling cloud event request"))
		return
	}

	logger.Info(fmt.Sprintf("Event TwinInstance: %s - Event TwinInterface: %s", twinEvent.TwinInstance, twinEvent.TwinInterface))

	err := ktwin.HandleEvent(twinEvent, model.TWIN_INTERFACE_PARKING_SPOT, handleParkingSpotEvent)

	if err != nil {
		logger.Error("Error processing cloud event request", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing cloud event request"))
		return
	}
}

func handleParkingSpotEvent(event *ktwin.TwinEvent) error {
	var parkingSpot model.ParkingSpot

	err := event.ToModel(&parkingSpot)

	if err != nil {
		return err
	}

	if parkingSpot.Status == "" {
		logger.Error(fmt.Sprintf("ParkingSpot status is empty for instance %s", event.TwinInstance), nil)
		return nil
	}

	latestParkingSpot, err := keventstore.GetLatestTwinEvent(event.TwinInstance, event.TwinInterface)

	if err != nil {
		return err
	}

	if latestParkingSpot == nil {
		logger.Info(fmt.Sprintf("No previous parking spot event found for instance %s", event.TwinInstance))
		return nil
	}

	if parkingSpot.Status == model.Occupied {
		updateCommand := parkingModel.UpdateVehicleCountCommand{
			VehicleEntranceCount: 1,
		}
		return kcommand.PublishCommand(model.TWIN_COMMAND_PARKING_UPDATE_VEHICLE_COUNT, updateCommand, model.TWIN_INTERFACE_OFF_STREET_PARKING_RELATIONSHIP, latestParkingSpot.TwinInstance, twinGraph)
	}

	if parkingSpot.Status == model.Free {
		updateCommand := parkingModel.UpdateVehicleCountCommand{
			VehicleExitCount: 1,
		}
		return kcommand.PublishCommand(model.TWIN_COMMAND_PARKING_UPDATE_VEHICLE_COUNT, updateCommand, model.TWIN_INTERFACE_OFF_STREET_PARKING_RELATIONSHIP, latestParkingSpot.TwinInstance, twinGraph)
	}

	logger.Info(fmt.Sprintf("ParkingSpot status is not recognized for instance %s", event.TwinInstance))

	return nil
}

func main() {
	server.LoadEnv()

	var err error
	twinGraph, err = ktwingraph.LoadTwinGraphByInstances([]string{model.TWIN_INTERFACE_PARKING_SPOT})

	if err != nil {
		logger.Error("Error loading twin graph", err)
		return
	}
	server.StartServer(requestHandler)
}
