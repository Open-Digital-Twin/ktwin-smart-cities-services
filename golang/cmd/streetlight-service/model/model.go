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
	Circuit              string     `json:"circuit,omitempty"`
	Status               LampStatus `json:"status,omitempty"`
	PowerState           PowerState `json:"powerState,omitempty"`
	DateLastLampChange   *time.Time `json:"dateLastLampChange,omitempty"`
	DateLastSwitchingOn  *time.Time `json:"dateLastSwitchingOn,omitempty"`
	DateLastSwitchingOff *time.Time `json:"dateLastSwitchingOff,omitempty"`
	ControllingMethod    string     `json:"controllingMethod,omitempty"`
	DateServiceStarted   *time.Time `json:"dateServiceStarted,omitempty"`
	Image                string     `json:"image,omitempty"`
	Annotations          string     `json:"annotations,omitempty"`
	LanternHeight        float64    `json:"lanternHeight,omitempty"`
	IlluminanceLevel     float64    `json:"illuminanceLevel,omitempty"`
	LocationCategory     string     `json:"locationCategory,omitempty"`
}
