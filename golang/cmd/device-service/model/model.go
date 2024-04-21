package model

import "time"

const (
	TWIN_INTERFACE_DEVICE = "ngsi-ld-city-device"
)

type Device struct {
	DataProvider         string     `json:"dataProvider,omitempty"`
	BatteryLevel         float64    `json:"batteryLevel,omitempty"`
	MeasurementFrequency int        `json:"measurementFrequency,omitempty"`
	Source               string     `json:"source,omitempty"`
	DateCreated          *time.Time `json:"dateCreated,omitempty"`
	DateObserved         *time.Time `json:"dateObserved,omitempty"`
	DateModified         *time.Time `json:"dateModified,omitempty"`
}
