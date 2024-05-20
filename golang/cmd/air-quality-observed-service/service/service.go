package service

import (
	"fmt"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/air-quality-observed-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	ktwingraph "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"

	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var (
	TWIN_INTERFACE_CITY_POLE            = "city-pole"
	TWIN_INTERFACE_AIR_QUALITY_OBSERVED = "ngsi-ld-city-airqualityobserved"

	// City Pole Update Air Quality Index Command
	TWIN_COMMAND_AIR_QUALITY_CITY_POLE_UPDATE_AIR_QUALITY_INDEX = "updateAirQualityIndex"
	TWIN_COMMAND_AIR_QUALITY_CITY_POLE_RELATIONSHIP_NAME        = "citypole"
)

var logger = log.NewLogger()
var twinGraph *ktwin.TwinGraph

func loadTwinGraph() error {
	if twinGraph == nil {
		var err error
		graph, err := ktwingraph.LoadTwinGraphByInstances([]string{TWIN_INTERFACE_AIR_QUALITY_OBSERVED})
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

	return kevent.HandleEvent(event, TWIN_INTERFACE_AIR_QUALITY_OBSERVED, handleAirQualityObservedEvent)
}

func handleAirQualityObservedEvent(event *ktwin.TwinEvent) error {
	var airQualityObserved model.AirQualityEvent

	err := event.ToModel(&airQualityObserved)

	if err != nil {
		return err
	}

	airQualityObserved.CalcCOAqiLevel()
	airQualityObserved.CalcPM10AqiLevel()
	airQualityObserved.CalcPM25AqiLevel()
	airQualityObserved.CalcSO2AqiLevel()
	airQualityObserved.CalcO3AqiLevel()

	event.SetData(airQualityObserved)
	err = keventstore.UpdateTwinEvent(event)

	if err != nil {
		return err
	}

	allLevels := []model.AQICategory{
		airQualityObserved.COAqiLevel,
		airQualityObserved.PM10AqiLevel,
		airQualityObserved.PM25AqiLevel,
		airQualityObserved.SO2AqiLevel,
		airQualityObserved.O3AqiLevel,
	}

	var updateAirQualityIndexCommand model.UpdateAirQualityIndexCommand
	updateAirQualityIndexCommand.SetAqiLevel(allLevels)

	if twinGraph == nil {
		logger.Error("Twin Graph not loaded", nil)
		return nil
	}

	err = kcommand.PublishCommand(TWIN_COMMAND_AIR_QUALITY_CITY_POLE_UPDATE_AIR_QUALITY_INDEX, updateAirQualityIndexCommand, TWIN_COMMAND_AIR_QUALITY_CITY_POLE_RELATIONSHIP_NAME, event.TwinInstance, *twinGraph)

	if err != nil {
		logger.Error(fmt.Sprintf("Error executing command %s in relation %s in TwinInstance %s\n", TWIN_COMMAND_AIR_QUALITY_CITY_POLE_UPDATE_AIR_QUALITY_INDEX, TWIN_COMMAND_AIR_QUALITY_CITY_POLE_RELATIONSHIP_NAME, event.TwinInstance), err)
		return err
	}

	return nil
}
