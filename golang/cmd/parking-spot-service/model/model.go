package model

import "time"

var (
	TWIN_INTERFACE_ON_STREET_PARKING_SPOT  = "ngsi-ld-city-onstreetparkingspot"
	TWIN_INTERFACE_OFF_STREET_PARKING_SPOT = "ngsi-ld-city-offstreetparkingspot"

	TWIN_INTERFACE_PARKING_SPOT                    = "ngsi-ld-city-parkingspot"
	TWIN_INTERFACE_OFF_STREET_PARKING_RELATIONSHIP = "refOffStreetParking"
	TWIN_COMMAND_PARKING_UPDATE_VEHICLE_COUNT      = "updateVehicleCount"
)

type Category string
type Status string

const (
	OffStreet Category = "offStreet"
	OnStreet  Category = "onStreet"

	Occupied Status = "occupied"
	Free     Status = "free"
	Closed   Status = "closed"
	Unknown  Status = "unknown"
)

type ParkingSpot struct {
	DateObserved float64   `json:"dateObserved"`
	Width        float64   `json:"width"`
	Length       float64   `json:"length"`
	TimeInstant  time.Time `json:"timeInstant"`
	Image        string    `json:"image"`
	Color        string    `json:"color"`
	Category     Category  `json:"category"`
	Status       Status    `json:"status"`
}
