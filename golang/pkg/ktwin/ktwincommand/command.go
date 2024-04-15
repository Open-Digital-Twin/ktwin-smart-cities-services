package ktwincommand

import (
	"fmt"
	"strings"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/ktwingraph"
)

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

// func HandleCommand(request *http.Request, twinInterface, command string, twinGraph ktwin.TwinGraph, callback func(*ktwin.KTwinCommandEvent, ktwin.TwinInstanceReference)) {
// 	ktwinCommandEvent := ktwin.HandleRequest(request)
// 	targetTwinInstance := ktwingraph.GetTwinGraphByRelation(ktwinCommandEvent.TwinInterface, ktwinCommandEvent.TwinInstance, twinGraph)

// 	if targetTwinInstance == nil {
// 		// TODO: need to handle the scenario where a TwinInterface has multiple relations with the same TwinInterface
// 		logger.Info(fmt.Sprintf("Twin Instance %s does not have a relation with the target interface: %s", ktwinCommandEvent.TwinInstance, ktwinCommandEvent.TwinInterface))
// 		return
// 	}

// 	if ktwinCommandEvent.TwinInterface == twinInterface && strings.ToLower(ktwinCommandEvent.Command) == strings.ToLower(command) {
// 		callback(ktwinCommandEvent, targetTwinInstance)
// 	}
// }
