package service

import (
	"fmt"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/pole-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	ktwingraph "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"

	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var (
	TWIN_INTERFACE_CITY_POLE = "city-pole"

	// City Pole Update Air Quality Index Command
	TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX = "updateAirQualityIndex"
	TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_RELATIONSHIP_NAME        = "refNeighborhood"
)

var logger = log.NewLogger()
var twinGraph *ktwin.TwinGraph

func loadTwinGraph() error {
	if twinGraph == nil {
		var err error
		graph, err := ktwingraph.LoadTwinGraphByInstances([]string{TWIN_INTERFACE_CITY_POLE})
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

	return kcommand.HandleCommand(event, TWIN_INTERFACE_CITY_POLE, TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX, *twinGraph, handleCityPoleCommand)
}

func handleCityPoleCommand(event *ktwin.TwinEvent) error {
	var updateAirQualityIndexCommand model.UpdateAirQualityIndexCommand
	err := event.ToModel(&updateAirQualityIndexCommand)

	if err != nil {
		return err
	}

	if twinGraph == nil {
		logger.Error("Twin Graph not loaded", nil)
		return nil
	}

	err = kcommand.PublishCommand(TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX, updateAirQualityIndexCommand, TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_RELATIONSHIP_NAME, event.TwinInstance, *twinGraph)

	if err != nil {
		logger.Error(fmt.Sprintf("Error executing command %s in relation %s in TwinInstance %s\n", TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX, TWIN_COMMAND_CITY_POLE_NEIGHBORHOOD_RELATIONSHIP_NAME, event.TwinInstance), err)
		return err
	}

	return nil
}
