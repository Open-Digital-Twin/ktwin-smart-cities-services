package model

import "time"

const (
	TWIN_INTERFACE_DEVICE = "ngsi-ld-city-device"
)

type Device struct {
	DataProvider         string    `json:"dataProvider"`
	BatteryLevel         float64   `json:"batteryLevel"`
	MeasurementFrequency int       `json:"measurementFrequency"`
	Source               string    `json:"source"`
	DateCreated          time.Time `json:"dateCreated"`
	DateObserved         time.Time `json:"dateObserved"`
	DateModified         time.Time `json:"dateModified"`
}
