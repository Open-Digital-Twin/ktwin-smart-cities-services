package service

import (
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/traffic-flow-observed-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
)

var (
	TWIN_INTERFACE_TRAFFIC_FLOW_OBSERVED = "ngsi-ld-city-trafficflowobserved"

	TRAFFIC_FLOW_AVERAGE_TRAFFIC_SPEED_THRESHOLD = 12
	TRAFFIC_FLOW_HEADWAY_TIME_THRESHOLD          = 2
)

func HandleEvent(event *ktwin.TwinEvent) error {
	return kevent.HandleEvent(event, TWIN_INTERFACE_TRAFFIC_FLOW_OBSERVED, handleTrafficFlowObservedEvent)
}

func handleTrafficFlowObservedEvent(event *ktwin.TwinEvent) error {
	var trafficFlowObserved model.TrafficFlowObservedEvent

	err := event.ToModel(&trafficFlowObserved)
	if err != nil {
		return err
	}

	if trafficFlowObserved.AverageVehicleSpeed < float64(TRAFFIC_FLOW_AVERAGE_TRAFFIC_SPEED_THRESHOLD) {
		trafficFlowObserved.Congested = true
	} else if trafficFlowObserved.AverageHeadwayTime < float64(TRAFFIC_FLOW_HEADWAY_TIME_THRESHOLD) {
		trafficFlowObserved.Congested = true
	} else {
		trafficFlowObserved.Congested = false
	}

	event.SetData(trafficFlowObserved)
	return keventstore.UpdateTwinEvent(event)
}
