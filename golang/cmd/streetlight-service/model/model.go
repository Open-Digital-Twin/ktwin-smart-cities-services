package model

import "time"

const (
	STREETLIGHT_INTERFACE_ID = "ngsi-ld-city-streetlight"
)

type PowerState string
type LampStatus string

const (
	PowerOn        PowerState = "on"
	PowerOff       PowerState = "off"
	PowerLow       PowerState = "low"
	PowerBootingUp PowerState = "bootingUp"
)

const (
	LampStatusOk            LampStatus = "oik"
	LampStatusDefective     LampStatus = "defectiveLamp"
	LampStatusColumnIssue   LampStatus = "columnIssue"
	LampStatusBrokenLantern LampStatus = "brokenLantern"
)

type Streetlight struct {
	Circuit              string     `json:"circuit"`
	Status               LampStatus `json:"status"`
	PowerState           PowerState `json:"powerState"`
	DateLastLampChange   time.Time  `json:"dateLastLampChange"`
	DateLastSwitchingOn  time.Time  `json:"dateLastSwitchingOn"`
	DateLastSwitchingOff time.Time  `json:"dateLastSwitchingOff"`
	ControllingMethod    string     `json:"controllingMethod"`
	DateServiceStarted   time.Time  `json:"dateServiceStarted"`
	Image                string     `json:"image"`
	Annotations          string     `json:"annotations"`
	LanternHeight        float64    `json:"lanternHeight"`
	IlluminanceLevel     float64    `json:"illuminanceLevel"`
	LocationCategory     string     `json:"locationCategory"`
}
