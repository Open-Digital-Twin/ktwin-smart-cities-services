package model

import "time"

type TrafficFlowObservedEvent struct {
	AreaServed           string     `json:"areaServed,omitempty"`           // The geographic area where a service or offered item is provided
	AverageGapDistance   float64    `json:"averageGapDistance,omitempty"`   // Average gap distance between consecutive vehicles
	AverageHeadwayTime   float64    `json:"averageHeadwayTime,omitempty"`   // Average headway time (time elapsed between two consecutive vehicles)
	AverageVehicleLength float64    `json:"averageVehicleLength,omitempty"` // Average length of the vehicles transiting during the observation period
	AverageVehicleSpeed  float64    `json:"averageVehicleSpeed,omitempty"`  // Average speed of the vehicles transiting during the observation period
	Congested            bool       `json:"congested"`                      // Flags whether there was a traffic congestion during the observation period in the referred lane
	DateObservedFrom     *time.Time `json:"dateObservedFrom,omitempty"`     // Observation period start date and time
	DateObservedTo       *time.Time `json:"dateObservedTo,omitempty"`       // Observation period end date and time
	Intensity            int        `json:"intensity,omitempty"`            // Total number of vehicles detected during this observation period
	LaneDirection        string     `json:"laneDirection,omitempty"`        // Usual direction of travel in the lane referred by this observation ("forward" or "backward")
	Occupancy            float64    `json:"occupancy,omitempty"`            // Fraction of the observation time where a vehicle has been occupying the observed lane
	ReversedLane         bool       `json:"reversedLane,omitempty"`         // Flags whether traffic in the lane was reversed during the observation period
	VehicleSubType       string     `json:"vehicleSubType,omitempty"`       // Subtype of vehicleType, provides more specific information about the type of vehicle
	VehicleType          string     `json:"vehicleType,omitempty"`          // Type of vehicle from the point of view of its structural characteristics
}
