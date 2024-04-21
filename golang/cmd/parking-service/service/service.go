package service

import (
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var logger = log.NewLogger()
var twinGraph *ktwin.TwinGraph

func loadTwinGraph() error {
	if twinGraph == nil {
		var err error
		graph, err := ktwingraph.LoadTwinGraphByInstances([]string{model.TWIN_INTERFACE_OFF_STREET_PARKING})
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
	return kcommand.HandleCommand(event, model.TWIN_INTERFACE_OFF_STREET_PARKING, model.TWIN_COMMAND_UPDATE_VEHICLE_COUNT, *twinGraph, handleUpdateVehicleCountCommand)
}

func handleUpdateVehicleCountCommand(command *ktwin.TwinEvent) error {
	var parking model.OffStreetParking
	var commandPayload model.UpdateVehicleCountCommand
	err := command.ToModel(&commandPayload)
	if err != nil {
		return err
	}

	if commandPayload.VehicleEntranceCount == 0 && commandPayload.VehicleExitCount == 0 {
		logger.Info("Vehicle entrance and exit count are 0, no need to update the twin")
		return nil
	}

	latestEvent, err := keventstore.GetLatestTwinEvent(command.TwinInterface, command.TwinInstance)

	if err != nil {
		return err
	}

	if latestEvent == nil {
		parking.TotalSpotNumber = 50 // default value
		parking.OccupiedSpotNumber = 0

		if commandPayload.VehicleEntranceCount != 0 {
			parking.IncrementOccupiedSpotNumber()
		} else {
			parking.DecrementOccupiedSpotNumber()
		}

		newEvent := ktwin.NewTwinEvent()
		newEvent.SetEvent(command.TwinInterface, command.TwinInstance, ktwin.RealEvent, parking)
		newEvent.SetData(parking)
		return keventstore.UpdateTwinEvent(newEvent)
	}

	err = latestEvent.ToModel(&parking)
	if err != nil {
		return err
	}

	if commandPayload.VehicleEntranceCount == 0 {
		logger.Info("Vehicle entrance count is 0, no need to update the twin")
	} else {
		parking.IncrementOccupiedSpotNumber()
	}

	if commandPayload.VehicleExitCount == 0 {
		logger.Info("Vehicle exit count is 0, no need to update the twin")
	} else {
		parking.DecrementOccupiedSpotNumber()
	}

	latestEvent.SetData(parking)
	return keventstore.UpdateTwinEvent(latestEvent)
}
