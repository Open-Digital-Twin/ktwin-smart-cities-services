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
	ceSource := relationship.Instance
	cloudEvent := ktwin.BuildCloudEvent(ceType, ceSource, commandPayload)

	fmt.Printf("Publishing Command Ce-Type: %s - Publishing Ce-Source: %s\n", ceType, ceSource)

	err := ktwin.PostCloudEvent(cloudEvent, ktwin.GetBrokerURL())

	if err != nil {
		return err
	}

	return nil
}

func HandleCommand(twinEvent *ktwin.TwinEvent, twinInterface string, command string, twinGraph ktwin.TwinGraph, callback func(*ktwin.TwinEvent) error) error {
	if twinEvent.EventType == ktwin.CommandEvent && strings.EqualFold(twinEvent.TwinInterface, twinInterface) {
		if strings.EqualFold(twinEvent.TwinInstance, twinEvent.TwinInstance) && strings.EqualFold(twinEvent.CommandName, command) {
			return callback(twinEvent)
		}
	}
	return nil
}
