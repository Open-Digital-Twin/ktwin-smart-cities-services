package ktwin

import (
	"fmt"
	"net/http"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

func PublishToRealTwin(twinInterface, twinInstance string, data interface{}) error {
	ceType := fmt.Sprintf(EventVirtualGenerated, twinInterface)
	ceSource := twinInstance
	cloudEvent := BuildCloudEvent(ceType, ceSource, data)
	return PostCloudEvent(cloudEvent, GetBrokerURL())
}

func PublishToVirtualTwin(twinInterface, twinInstance string, data interface{}) error {
	ceType := fmt.Sprintf(EventRealGenerated, twinInterface)
	ceSource := twinInstance
	cloudEvent := BuildCloudEvent(ceType, ceSource, data)
	return PostCloudEvent(cloudEvent, GetBrokerURL())
}

func HandleRequest(r *http.Request) *TwinEvent {
	// Handle incoming events
	twinEvent := NewTwinEvent()
	err := twinEvent.HandleRequest(r)

	if err != nil {
		logger := logger.NewLogger()
		logger.Error("Error handling cloud event request", err)
		return nil
	}

	return twinEvent
}

func HandleEvent(twinEvent *TwinEvent, twinInterface string, callback func(*TwinEvent) error) error {
	if twinEvent.TwinInterface == twinInterface {
		return callback(twinEvent)
	}
	return nil
}
