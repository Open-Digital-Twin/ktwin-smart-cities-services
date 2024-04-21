package service

import (
	"fmt"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/pole-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kcommand"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	ktwingraph "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"

	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var (
	TWIN_INTERFACE_AIR_QUALITY_OBSERVED  = "ngsi-ld-city-airqualityobserved"
	TWIN_INTERFACE_WEATHER_OBSERVED      = "ngsi-ld-city-weatherobserved"
	TWIN_INTERFACE_CROWD_FLOW_OBSERVED   = "ngsi-ld-city-crowdflowobserved"
	TWIN_INTERFACE_TRAFFIC_FLOW_OBSERVED = "ngsi-ld-city-trafficflowobserved"

	TWIN_COMMAND_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX    = "updateAirQualityIndex"
	TWIN_COMMAND_RELATIONSHIP_NEIGHBORHOOD_UPDATE_WEATHER = "refNeighborhood"

	CROWD_FLOW_AVERAGE_CROWD_SPEED_THRESHOLD = 4
	CROWD_FLOW_HEADWAY_TIME_THRESHOLD        = 2

	TRAFFIC_FLOW_AVERAGE_TRAFFIC_SPEED_THRESHOLD = 12
	TRAFFIC_FLOW_HEADWAY_TIME_THRESHOLD          = 2
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

	switch event.TwinInterface {
	case TWIN_INTERFACE_AIR_QUALITY_OBSERVED:
		return handleAirQualityObservedEvent(event)
	case TWIN_INTERFACE_CROWD_FLOW_OBSERVED:
		return handleCrowdFlowObservedEvent(event)
	case TWIN_INTERFACE_TRAFFIC_FLOW_OBSERVED:
		return handleTrafficFlowObservedEvent(event)
	case TWIN_INTERFACE_WEATHER_OBSERVED:
		return handleWeatherObservedEvent(event)
	default:
		logger.Info(fmt.Sprintf("Unhandled event for interface: %s\n", event.TwinInterface))
	}
	return nil
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
	for _, level := range allLevels {
		if level == model.Hazardous {
			updateAirQualityIndexCommand.AqiLevel = model.Hazardous
			break
		}
		if level == model.VeryUnhealthy {
			updateAirQualityIndexCommand.AqiLevel = model.VeryUnhealthy
			break
		}
		if level == model.Unhealthy {
			updateAirQualityIndexCommand.AqiLevel = model.Unhealthy
			break
		}
		if level == model.UnhealthyForSensitiveGroups {
			updateAirQualityIndexCommand.AqiLevel = model.UnhealthyForSensitiveGroups
			break
		}
		if level == model.Moderate {
			updateAirQualityIndexCommand.AqiLevel = model.Moderate
			break
		}
	}

	cityPoleRelation := ktwingraph.GetRelationshipFromGraph(event.TwinInstance, "citypole", *twinGraph)
	if cityPoleRelation == nil {
		logger.Error(fmt.Sprintf("City pole relation not found for Twin Instance: %s\n", event.TwinInstance), nil)
	}

	if twinGraph == nil {
		logger.Error("Twin Graph not loaded", nil)
		return nil
	}

	err = kcommand.PublishCommand(TWIN_COMMAND_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX, updateAirQualityIndexCommand, TWIN_COMMAND_RELATIONSHIP_NEIGHBORHOOD_UPDATE_WEATHER, cityPoleRelation.Instance, *twinGraph)

	if err != nil {
		logger.Error(fmt.Sprintf("Error executing command %s in relation %s in TwinInstance %s\n", TWIN_COMMAND_NEIGHBORHOOD_UPDATE_AIR_QUALITY_INDEX, TWIN_COMMAND_RELATIONSHIP_NEIGHBORHOOD_UPDATE_WEATHER, cityPoleRelation.Instance), err)
		return err
	}

	return nil
}

func handleCrowdFlowObservedEvent(event *ktwin.TwinEvent) error {
	var crowdFlowObserved model.CrowdFlowObservedEvent

	err := event.ToModel(&crowdFlowObserved)
	if err != nil {
		return err
	}

	if crowdFlowObserved.AverageCrowdSpeed < float64(CROWD_FLOW_AVERAGE_CROWD_SPEED_THRESHOLD) {
		crowdFlowObserved.Congested = true
	} else if crowdFlowObserved.AverageHeadwayTime < float64(CROWD_FLOW_HEADWAY_TIME_THRESHOLD) {
		crowdFlowObserved.Congested = true
	} else {
		crowdFlowObserved.Congested = false
	}

	event.SetData(crowdFlowObserved)

	return keventstore.UpdateTwinEvent(event)
}

func handleTrafficFlowObservedEvent(event *ktwin.TwinEvent) error {
	var trafficFlowObserved model.TrafficFlowObservedEvent

	err := event.ToModel(&trafficFlowObserved)
	if err != nil {
		return err
	}

	if trafficFlowObserved.AverageVehicleSpeed < float64(TRAFFIC_FLOW_AVERAGE_TRAFFIC_SPEED_THRESHOLD) {
		trafficFlowObserved.Congested = true
	} else if trafficFlowObserved.AverageHeadwayTime < float64(TRAFFIC_FLOW_HEADWAY_TIME_THRESHOLD) {
		trafficFlowObserved.Congested = true
	} else {
		trafficFlowObserved.Congested = false
	}

	event.SetData(trafficFlowObserved)
	return keventstore.UpdateTwinEvent(event)
}

func handleWeatherObservedEvent(event *ktwin.TwinEvent) error {
	latestEvent, err := keventstore.GetLatestTwinEvent(event.TwinInterface, event.TwinInstance)

	if err != nil {
		return err
	}

	if latestEvent == nil {
		latestEvent = event
	}

	var latestWeatherObserved model.WeatherObservedEvent
	err = latestEvent.ToModel(&latestWeatherObserved)

	if err != nil {
		return err
	}

	var weatherObserved model.WeatherObservedEvent
	err = event.ToModel(&weatherObserved)

	if err != nil {
		return err
	}

	weatherObserved.SetPressureTendency(latestWeatherObserved.AtmosphericPressure)
	weatherObserved.SetFeelsLikeTemperature(weatherObserved.Temperature, weatherObserved.WindSpeed)
	weatherObserved.SetDewpoint(weatherObserved.Temperature, weatherObserved.RelativeHumidity)

	event.SetData(weatherObserved)
	return keventstore.UpdateTwinEvent(event)
}
