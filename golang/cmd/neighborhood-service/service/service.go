package service

import (
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/neighborhood-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/clock"
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
		graph, err := ktwingraph.LoadTwinGraphByInstances([]string{model.TWIN_INTERFACE_NEIGHBORHOOD, model.TWIN_INTERFACE_CITY_POLE})
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
	return kcommand.HandleCommand(event, model.TWIN_COMMAND_UPDATE_AIR_QUALITY_INDEX, *twinGraph, handleUpdateAirQualityIndex)
}

func handleUpdateAirQualityIndex(command *ktwin.TwinEvent, targetTwinInstance ktwin.TwinInstanceReference) error {
	latestEvent, err := keventstore.GetLatestTwinEvent(targetTwinInstance.Interface, targetTwinInstance.Instance)

	if err != nil {
		return err
	}

	var neighborhood model.Neighborhood

	if latestEvent == nil {
		now := clock.Now()
		neighborhood = model.Neighborhood{
			AqiLevel:     model.GOOD,
			DateObserved: now,
			DateModified: now,
		}
		latestEvent = ktwin.NewTwinEvent()
		latestEvent.SetEvent(model.TWIN_INTERFACE_NEIGHBORHOOD, targetTwinInstance.Instance, ktwin.RealEvent, neighborhood)
	} else {
		err = latestEvent.ToModel(&neighborhood)
		if err != nil {
			return err
		}
	}

	var updateAirQualityIndexCommand model.UpdateAirQualityIndexCommand
	err = command.ToModel(&updateAirQualityIndexCommand)
	if err != nil {
		return err
	}

	if updateAirQualityIndexCommand.AqiLevel == "" {
		logger.Info("AqiLevel not provided")
		return nil
	}

	newQualityIndexInt := model.GetQualityLevelInteger(updateAirQualityIndexCommand.AqiLevel)
	latestQualityIndexInt := model.GetQualityLevelInteger(neighborhood.AqiLevel)

	if newQualityIndexInt > latestQualityIndexInt || hasTimeExpired(clock.Now(), neighborhood.DateObserved, 60) {
		neighborhood.AqiLevel = updateAirQualityIndexCommand.AqiLevel
		neighborhood.DateModified = clock.Now()
	}

	latestEvent.SetData(neighborhood)

	return keventstore.UpdateTwinEvent(latestEvent)
}

func hasTimeExpired(datetimeNow time.Time, datetimeObserved time.Time, minutes int) bool {
	return datetimeNow.Sub(datetimeObserved).Minutes() > float64(minutes)
}
