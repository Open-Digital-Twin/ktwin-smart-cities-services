package kevent

import (
	"fmt"
	"net/http"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var logger = log.NewLogger()

func PublishToRealTwin(twinInterface, twinInstance string, data interface{}) error {
	ceType := fmt.Sprintf(ktwin.EventVirtualGenerated, twinInterface)
	ceSource := twinInstance
	cloudEvent := ktwin.BuildCloudEvent(ceType, ceSource, data)
	return ktwin.PostCloudEvent(cloudEvent, ktwin.GetBrokerURL())
}

func PublishToVirtualTwin(twinInterface, twinInstance string, data interface{}) error {
	ceType := fmt.Sprintf(ktwin.EventRealGenerated, twinInterface)
	ceSource := twinInstance
	cloudEvent := ktwin.BuildCloudEvent(ceType, ceSource, data)
	return ktwin.PostCloudEvent(cloudEvent, ktwin.GetBrokerURL())
}

func HandleRequest(r *http.Request) *ktwin.TwinEvent {
	// Handle incoming events
	twinEvent := ktwin.NewTwinEvent()
	err := twinEvent.HandleRequest(r)

	if err != nil {
		logger.Error("Error handling cloud event request", err)
		return nil
	}

	return twinEvent
}

func HandleEvent(twinEvent *ktwin.TwinEvent, twinInterface string, callback func(*ktwin.TwinEvent) error) error {
	if twinEvent.TwinInterface == twinInterface {
		return callback(twinEvent)
	}
	return nil
}

func RequestHandlerFunc(w http.ResponseWriter, r *http.Request, handleEvent func(*ktwin.TwinEvent) error) {
	twinEvent := HandleRequest(r)

	if twinEvent == nil {
		logger.Error("Error handling cloud event request", nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error handling cloud event request"))
		return
	}

	if err := handleEvent(twinEvent); err != nil {
		logger.Error("Error processing cloud event request", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing cloud event request"))
		return
	}
}
