package service

import (
	"fmt"

	parkingModel "github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-spot-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var logger = log.NewLogger()
var twinGraph *ktwin.TwinGraph

func loadTwinGraph() error {
	if twinGraph == nil {
		var err error
		graph, err := ktwingraph.LoadTwinGraphByInterfaces([]string{model.TWIN_INTERFACE_PARKING_SPOT})
		if err != nil {
			logger.Error("Error loading twin graph", err)
			return err
		}
		twinGraph = &graph
	}
	return nil
}

func HandleEvent(event *ktwin.TwinEvent) error {
	err := loadTwinGraph()
	if err != nil {
		return err
	}
	return kevent.HandleEvent(event, model.TWIN_INTERFACE_PARKING_SPOT, handleParkingSpotEvent)
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

	if parkingSpot.Status == model.Occupied {
		updateCommand := parkingModel.UpdateVehicleCountCommand{
			VehicleEntranceCount: 1,
		}
		return kcommand.PublishCommand(model.TWIN_COMMAND_PARKING_UPDATE_VEHICLE_COUNT, updateCommand, model.TWIN_INTERFACE_OFF_STREET_PARKING_RELATIONSHIP, event.TwinInstance, *twinGraph)
	}

	if parkingSpot.Status == model.Free {
		updateCommand := parkingModel.UpdateVehicleCountCommand{
			VehicleExitCount: 1,
		}
		return kcommand.PublishCommand(model.TWIN_COMMAND_PARKING_UPDATE_VEHICLE_COUNT, updateCommand, model.TWIN_INTERFACE_OFF_STREET_PARKING_RELATIONSHIP, event.TwinInstance, *twinGraph)
	}

	logger.Info(fmt.Sprintf("ParkingSpot status is not recognized for instance %s", event.TwinInstance))

	return nil
}
