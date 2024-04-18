package model

import "time"

var (
	TWIN_INTERFACE_NEIGHBORHOOD           = "s4city-city-neighborhood"
	TWIN_INTERFACE_CITY_POLE              = "city-pole"
	TWIN_COMMAND_UPDATE_AIR_QUALITY_INDEX = "updateAirQualityIndex"
)

type AQICategory string

func GetQualityLevelInteger(aqi AQICategory) int {
	switch aqi {
	case GOOD:
		return 1
	case MODERATE:
		return 2
	case UNHEALTHY_FOR_SENSITIVE_GROUPS:
		return 3
	case UNHEALTHY:
		return 4
	case VERY_UNHEALTHY:
		return 5
	case HAZARDOUS:
		return 6
	default:
		return 0
	}
}

const (
	GOOD                           AQICategory = "GOOD"
	MODERATE                       AQICategory = "MODERATE"
	UNHEALTHY_FOR_SENSITIVE_GROUPS AQICategory = "UNHEALTHY_FOR_SENSITIVE_GROUPS"
	UNHEALTHY                      AQICategory = "UNHEALTHY"
	VERY_UNHEALTHY                 AQICategory = "VERY_UNHEALTHY"
	HAZARDOUS                      AQICategory = "HAZARDOUS"
)

type UpdateAirQualityIndexCommand struct {
	AqiLevel AQICategory `json:"aqiLevel"`
}

type Neighborhood struct {
	AqiLevel     AQICategory `json:"aqiLevel"`
	DateObserved time.Time   `json:"dateObserved"`
	DateModified time.Time   `json:"dateModified"`
}
