package model

import "time"

type CrowdFlowObservedEvent struct {
	DateObservedFrom   *time.Time `json:"dateObservedFrom,omitempty"`   // Date and time when the observation started
	DateObservedTo     *time.Time `json:"dateObservedTo,omitempty"`     // Date and time when the observation ended
	Occupancy          bool       `json:"occupancy,omitempty"`          // Fraction of the observation time where a person has been occupying the observed walkway
	AverageCrowdSpeed  float64    `json:"averageCrowdSpeed,omitempty"`  // Average speed of the crowd transiting during the observation period (Km/h)
	Congested          bool       `json:"congested"`                    // Flags whether there was a crowd congestion during the observation period in the referred walkway
	AverageHeadwayTime float64    `json:"averageHeadwayTime,omitempty"` // Average headway time (time elapsed between two consecutive persons) in seconds
	Direction          string     `json:"direction,omitempty"`          // Usual direction of travel in the walkway referred by this observation with respect to the city center ("inbound" or "outbound")
	PeopleCount        int        `json:"peopleCount,omitempty"`        // Number of people observed
}
