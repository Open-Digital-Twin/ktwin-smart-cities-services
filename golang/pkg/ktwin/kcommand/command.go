package kcommand

import (
	"fmt"
	"strings"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
)

// TwinCommand

func PublishCommand(command string, commandPayload interface{}, relationshipName, twinInstanceSource string, twinGraph ktwin.TwinGraph) error {
	relationship := ktwingraph.GetRelationshipFromGraph(twinInstanceSource, relationshipName, twinGraph)
	if relationship == nil {
		return fmt.Errorf("relationship %s not found in Twin Instance %s", relationshipName, twinInstanceSource)
	}
	ceType := fmt.Sprintf(ktwin.EventCommandExecuted, relationship.Interface, strings.ToLower(command))
	ceSource := twinInstanceSource
	cloudEvent := ktwin.BuildCloudEvent(ceType, ceSource, commandPayload)

	err := ktwin.PostCloudEvent(cloudEvent, ktwin.GetBrokerURL())

	if err != nil {
		return err
	}

	return nil
}

func HandleCommand(twinEvent *ktwin.TwinEvent, command string, twinGraph ktwin.TwinGraph, callback func(*ktwin.TwinEvent, ktwin.TwinInstanceReference) error) error {
	if twinEvent.EventType != ktwin.CommandEvent {
		targetTwinInstance := ktwingraph.GetTwinGraphByRelation(twinEvent.TwinInterface, twinEvent.TwinInstance, twinGraph)

		if targetTwinInstance == nil {
			// TODO: need to handle the scenario where a TwinInterface has multiple relations with the same TwinInterface
			return fmt.Errorf(fmt.Sprintf("Twin Instance %s does not have a relation with the target interface: %s", twinEvent.TwinInstance, twinEvent.TwinInterface))
		}

		if strings.EqualFold(twinEvent.TwinInstance, targetTwinInstance.Instance) && strings.EqualFold(twinEvent.CommandName, command) {
			return callback(twinEvent, *targetTwinInstance)
		}
	} else {
		return fmt.Errorf("event is not a command event")
	}
	return nil
}
