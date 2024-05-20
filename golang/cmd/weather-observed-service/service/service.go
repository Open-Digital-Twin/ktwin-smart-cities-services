package service

import (
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/weather-observed-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
)

var (
	TWIN_INTERFACE_WEATHER_OBSERVED = "ngsi-ld-city-weatherobserved"
)

func HandleEvent(event *ktwin.TwinEvent) error {
	return kevent.HandleEvent(event, TWIN_INTERFACE_WEATHER_OBSERVED, handleWeatherObservedEvent)
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
