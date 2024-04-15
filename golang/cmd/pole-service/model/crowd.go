package model

import "time"

type CrowdFlowObservedEvent struct {
	DateObservedFrom   time.Time // Date and time when the observation started
	DateObservedTo     time.Time // Date and time when the observation ended
	Occupancy          bool      // Fraction of the observation time where a person has been occupying the observed walkway
	AverageCrowdSpeed  int       // Average speed of the crowd transiting during the observation period (Km/h)
	Congested          bool      // Flags whether there was a crowd congestion during the observation period in the referred walkway
	AverageHeadwayTime float64   // Average headway time (time elapsed between two consecutive persons) in seconds
	Direction          string    // Usual direction of travel in the walkway referred by this observation with respect to the city center ("inbound" or "outbound")
	PeopleCount        int       // Number of people observed
}
