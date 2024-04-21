package model

import "math"

const (
	WEATHER_OBSERVED_TWIN_INTERFACE = "ngsi-ld-city-airqualityobserved"
)

// Weather Observed Event

type PressureTendency string

const (
	Raising PressureTendency = "raising"
	Steady  PressureTendency = "steady"
	Falling PressureTendency = "falling"
)

// WeatherObservedEvent represents the structure for an weather event
type WeatherObservedEvent struct {
	WeatherType          string           `json:"weatherType,omitempty"`          // Type of weather
	StationCode          string           `json:"stationCode,omitempty"`          // Code of the weather station
	StationName          string           `json:"stationName,omitempty"`          // Name of the weather station
	PressureTendency     PressureTendency `json:"pressureTendency,omitempty"`     // Pressure tendency (Raising, Steady, Falling)
	AtmosphericPressure  float64          `json:"atmosphericPressure,omitempty"`  // Atmospheric pressure in hPa
	Dewpoint             float64          `json:"dewpoint,omitempty"`             // Dewpoint temperature in Celsius
	FeelsLikeTemperature float64          `json:"feelsLikeTemperature,omitempty"` // Feels like temperature in Celsius
	Temperature          float64          `json:"temperature,omitempty"`          // Temperature in Celsius
	Illuminance          float64          `json:"illuminance,omitempty"`          // Illuminance in lux
	Precipitation        float64          `json:"precipitation,omitempty"`        // Precipitation in mm
	RelativeHumidity     float64          `json:"relativeHumidity,omitempty"`     // Relative humidity in percentage
	SnowHeight           float64          `json:"snowHeight,omitempty"`           // Snow height in cm
	SolarRadiation       float64          `json:"solarRadiation,omitempty"`       // Solar radiation in W/m^2
	StreamGauge          float64          `json:"streamGauge,omitempty"`          // Stream gauge in m
	UVIndexMax           float64          `json:"uvIndexMax,omitempty"`           // Maximum UV index
	Visibility           float64          `json:"visibility,omitempty"`           // Visibility in km
	WindDirection        float64          `json:"windDirection,omitempty"`        // Wind direction in degrees
	WindSpeed            float64          `json:"windSpeed,omitempty"`            // Wind speed in m/s
}

func (w *WeatherObservedEvent) SetPressureTendency(latestAtmosphericPressure float64) {
	difference := w.AtmosphericPressure - latestAtmosphericPressure
	if difference < 0 {
		w.PressureTendency = "falling"
		return
	} else if difference > 0 {
		w.PressureTendency = "raising"
		return
	} else {
		w.PressureTendency = "steady"
	}
}

func (w *WeatherObservedEvent) SetFeelsLikeTemperature(temperature float64, windSpeed float64) {
	w.FeelsLikeTemperature = 33 + (10*math.Sqrt(windSpeed)+10.45-windSpeed)*(temperature-33)/22
}

func (w *WeatherObservedEvent) SetDewpoint(temperature float64, relativeHumidity float64) {
	w.Dewpoint = temperature - ((100 - relativeHumidity) / 5)
}
