package service

import (
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/crowd-flow-observed-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
)

var (
	TWIN_INTERFACE_CROWD_FLOW_OBSERVED = "ngsi-ld-city-crowdflowobserved"

	CROWD_FLOW_AVERAGE_CROWD_SPEED_THRESHOLD = 4
	CROWD_FLOW_HEADWAY_TIME_THRESHOLD        = 2
)

func HandleEvent(event *ktwin.TwinEvent) error {
	return kevent.HandleEvent(event, TWIN_INTERFACE_CROWD_FLOW_OBSERVED, handleCrowdFlowObservedEvent)
}

func handleCrowdFlowObservedEvent(event *ktwin.TwinEvent) error {
	var crowdFlowObserved model.CrowdFlowObservedEvent

	err := event.ToModel(&crowdFlowObserved)
	if err != nil {
		return err
	}

	if crowdFlowObserved.AverageCrowdSpeed < float64(CROWD_FLOW_AVERAGE_CROWD_SPEED_THRESHOLD) {
		crowdFlowObserved.Congested = true
	} else if crowdFlowObserved.AverageHeadwayTime < float64(CROWD_FLOW_HEADWAY_TIME_THRESHOLD) {
		crowdFlowObserved.Congested = true
	} else {
		crowdFlowObserved.Congested = false
	}

	event.SetData(crowdFlowObserved)

	return keventstore.UpdateTwinEvent(event)
}
