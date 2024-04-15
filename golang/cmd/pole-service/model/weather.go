package model

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
	WeatherType          string           // Type of weather
	StationCode          string           // Code of the weather station
	StationName          string           // Name of the weather station
	PressureTendency     PressureTendency // Pressure tendency (Raising, Steady, Falling)
	AtmosphericPressure  float64          // Atmospheric pressure in hPa
	Dewpoint             float64          // Dewpoint temperature in Celsius
	FeelsLikeTemperature float64          // Feels like temperature in Celsius
	Temperature          float64          // Temperature in Celsius
	Illuminance          float64          // Illuminance in lux
	Precipitation        float64          // Precipitation in mm
	RelativeHumidity     float64          // Relative humidity in percentage
	SnowHeight           float64          // Snow height in cm
	SolarRadiation       float64          // Solar radiation in W/m^2
	StreamGauge          float64          // Stream gauge in m
	UVIndexMax           float64          // Maximum UV index
	Visibility           float64          // Visibility in km
	WindDirection        float64          // Wind direction in degrees
	WindSpeed            float64          // Wind speed in m/s
}
