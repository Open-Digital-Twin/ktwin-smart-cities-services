package model

import "time"

type TrafficFlowObservedEvent struct {
	AreaServed           string    // The geographic area where a service or offered item is provided
	AverageGapDistance   float64   // Average gap distance between consecutive vehicles
	AverageHeadwayTime   float64   // Average headway time (time elapsed between two consecutive vehicles)
	AverageVehicleLength float64   // Average length of the vehicles transiting during the observation period
	AverageVehicleSpeed  float64   // Average speed of the vehicles transiting during the observation period
	Congested            bool      // Flags whether there was a traffic congestion during the observation period in the referred lane
	DateObservedFrom     time.Time // Observation period start date and time
	DateObservedTo       time.Time // Observation period end date and time
	Intensity            int       // Total number of vehicles detected during this observation period
	LaneDirection        string    // Usual direction of travel in the lane referred by this observation ("forward" or "backward")
	Occupancy            float64   // Fraction of the observation time where a vehicle has been occupying the observed lane
	ReversedLane         bool      // Flags whether traffic in the lane was reversed during the observation period
	VehicleSubType       string    // Subtype of vehicleType, provides more specific information about the type of vehicle
	VehicleType          string    // Type of vehicle from the point of view of its structural characteristics
}
