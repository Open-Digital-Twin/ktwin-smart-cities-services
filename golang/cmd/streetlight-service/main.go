package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/streetlight-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/server"
)

var logger = log.NewLogger()

func requestHandler(w http.ResponseWriter, r *http.Request) {
	twinEvent := ktwin.HandleRequest(r)

	if twinEvent == nil {
		logger.Error("Error handling cloud event request", nil)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error handling cloud event request"))
		return
	}

	logger.Info(fmt.Sprintf("Event TwinInstance: %s - Event TwinInterface: %s", twinEvent.TwinInstance, twinEvent.TwinInterface))

	err := ktwin.HandleEvent(twinEvent, model.STREETLIGHT_INTERFACE_ID, handleStreetLightEvent)

	if err != nil {
		logger.Error("Error processing cloud event request", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing cloud event request"))
		return
	}
}

func handleStreetLightEvent(event *ktwin.TwinEvent) error {
	timeNow := time.Now()

	var currentStreetlight model.Streetlight
	err := event.ToModel(&currentStreetlight)

	if err != nil {
		logger.Error("Error parsing event data", err)
		return nil
	}

	if currentStreetlight.PowerState == "" {
		logger.Info("Streetlight event has no powerState attribute value")
		return nil
	}

	latestEvent, err := keventstore.GetLatestTwinEvent(event.TwinInstance, event.TwinInterface)

	if err != nil {
		return err
	}

	if latestEvent == nil {
		if currentStreetlight.PowerState == model.PowerOn {
			currentStreetlight.DateLastSwitchingOn = timeNow
		}
		if currentStreetlight.PowerState == model.PowerOff {
			currentStreetlight.DateLastSwitchingOff = timeNow
		}
		return keventstore.UpdateTwinEvent(event)
	}

	var latestStreetlight model.Streetlight
	err = latestEvent.ToModel(&latestStreetlight)

	if err != nil {
		return err
	}

	if latestStreetlight.PowerState == currentStreetlight.PowerState {
		if currentStreetlight.PowerState == model.PowerOn {
			if isWithDefect(timeNow, currentStreetlight.DateLastSwitchingOn) {
				currentStreetlight.Status = model.LampStatusDefective
			}
			currentStreetlight.DateLastSwitchingOn = timeNow
		}
		if currentStreetlight.PowerState == model.PowerOff {
			if isWithDefect(timeNow, currentStreetlight.DateLastSwitchingOff) {
				currentStreetlight.Status = model.LampStatusDefective
			}
			currentStreetlight.DateLastSwitchingOff = timeNow
		}
	}

	return keventstore.UpdateTwinEvent(event)
}

// In case of 48h of no change in the state, we consider that lamp with a defect
func isWithDefect(datetimeNow time.Time, dateLastSwitching time.Time) bool {
	timeDifference := datetimeNow.Sub(dateLastSwitching)
	return timeDifference > time.Hour*48
}

func main() {
	server.LoadEnv()
	server.StartServer(requestHandler)
}
