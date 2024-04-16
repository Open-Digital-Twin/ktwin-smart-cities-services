package kcommand

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// TwinCommand

var logger = log.NewLogger()

type TwinCommand struct {
	CloudEvent         *cloudevents.Event
	TwinInterface      string
	Command            string
	TwinInstanceSource string
}

func NewTwinCommandEvent() *TwinCommand {
	return &TwinCommand{}
}

func (c *TwinCommand) HandleRequest(r *http.Request) error {
	cloudEvent, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		logger.Error("failed to parse CloudEvent from request: %v", err)
		return err
	}
	parts := strings.Split(cloudEvent.Type(), ".")
	if len(parts) > 2 {
		c.TwinInterface = parts[2]
	}
	if len(parts) > 3 {
		c.Command = parts[3]
	}
	c.TwinInstanceSource = cloudEvent.Source()
	c.CloudEvent = cloudEvent
	return nil
}

func (k *TwinCommand) ToModel(model interface{}) error {
	err := json.Unmarshal(k.CloudEvent.Data(), model)
	if err != nil {
		return err
	}
	return nil
}

func PublishCommand(command string, commandPayload interface{}, relationshipName, twinInstanceSource string, twinGraph ktwin.TwinGraph) error {
	relationship := ktwingraph.GetRelationshipFromGraph(twinInstanceSource, relationshipName, twinGraph)
	if relationship == nil {
		return fmt.Errorf("relationship %s not found in Twin Instance %s", relationshipName, twinInstanceSource)
	}
	ceType := fmt.Sprintf(ktwin.EventCommandExecuted, relationship.Interface+"."+strings.ToLower(command))
	ceSource := twinInstanceSource
	cloudEvent := ktwin.BuildCloudEvent(ceType, ceSource, commandPayload)

	err := ktwin.PostCloudEvent(cloudEvent, ktwin.GetBrokerURL())

	if err != nil {
		return err
	}

	return nil
}

func HandleCommand(twinCommand *TwinCommand, command string, twinGraph ktwin.TwinGraph, callback func(*TwinCommand, ktwin.TwinInstanceReference) error) error {
	targetTwinInstance := ktwingraph.GetTwinGraphByRelation(twinCommand.TwinInterface, twinCommand.TwinInstanceSource, twinGraph)

	if targetTwinInstance == nil {
		// TODO: need to handle the scenario where a TwinInterface has multiple relations with the same TwinInterface
		return fmt.Errorf(fmt.Sprintf("Twin Instance %s does not have a relation with the target interface: %s", twinCommand.TwinInstanceSource, twinCommand.TwinInterface))
	}

	if strings.EqualFold(twinCommand.TwinInstanceSource, targetTwinInstance.Instance) && strings.EqualFold(twinCommand.Command, command) {
		return callback(twinCommand, *targetTwinInstance)
	}

	return nil
}

func HandleCommandRequest(r *http.Request) *TwinCommand {
	// Handle incoming events
	twinEvent := NewTwinCommandEvent()
	err := twinEvent.HandleRequest(r)

	if err != nil {
		logger.Error("Error handling cloud event request", err)
		return nil
	}

	return twinEvent
}
