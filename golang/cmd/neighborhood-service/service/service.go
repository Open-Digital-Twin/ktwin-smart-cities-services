package service

import (
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var (
	TWIN_INTERFACE_OFF_STREET_PARKING = "ngsi-ld-city-offstreetparking"
	TWIN_COMMAND_UPDATE_VEHICLE_COUNT = "updateVehicleCount"
)

var logger = log.NewLogger()
var twinGraph *ktwin.TwinGraph

func loadTwinGraph() error {
	if twinGraph != nil {
		var err error
		graph, err := ktwingraph.LoadTwinGraphByInstances([]string{model.TWIN_COMMAND_UPDATE_VEHICLE_COUNT})
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
	return kcommand.HandleCommand(event, TWIN_COMMAND_UPDATE_VEHICLE_COUNT, *twinGraph, handleUpdateVehicleCountCommand)
}

func handleUpdateVehicleCountCommand(command *ktwin.TwinEvent, targetTwinInstance ktwin.TwinInstanceReference) error {
	latestEvent, err := keventstore.GetLatestTwinEvent(targetTwinInstance.Instance, targetTwinInstance.Instance)

	if err != nil {
		return err
	}

	var parking model.OffStreetParking

	if latestEvent == nil {
		parking.OccupiedSpotNumber = 0
		// ktwin.NewTwinEvent().SetEvent(command.TwinInterface, command.TwinInstanceSource, command.TwinInstanceSource, parking)
	}

	var commandPayload model.UpdateVehicleCountCommand
	err = command.ToModel(&commandPayload)
	if err != nil {
		return err
	}

	if commandPayload.VehicleEntranceCount == 0 {
		logger.Info("Vehicle entrance count is 0, no need to update the twin")
	} else {
		parking.OccupiedSpotNumber = parking.OccupiedSpotNumber + commandPayload.VehicleEntranceCount
	}

	if commandPayload.VehicleExitCount == 0 {
		logger.Info("Vehicle exit count is 0, no need to update the twin")
	} else {
		if parking.OccupiedSpotNumber <= 0 {
			logger.Info("Vehicle exit count is greater than occupied spot number, no need to update the twin")
			return nil
		}
		parking.OccupiedSpotNumber = parking.OccupiedSpotNumber - commandPayload.VehicleExitCount
	}

	latestEvent.SetData(parking)
	return keventstore.UpdateTwinEvent(latestEvent)
}
