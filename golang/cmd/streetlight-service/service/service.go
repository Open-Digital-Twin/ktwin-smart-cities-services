package service

import (
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/streetlight-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/clock"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var logger = log.NewLogger()

func HandleEvent(event *ktwin.TwinEvent) error {
	return kevent.HandleEvent(event, model.STREETLIGHT_INTERFACE_ID, handleStreetLightEvent)
}

func handleStreetLightEvent(event *ktwin.TwinEvent) error {
	timeNow := clock.Now()

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

	latestEvent, err := keventstore.GetLatestTwinEvent(event.TwinInterface, event.TwinInstance)

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
		event.SetData(currentStreetlight)
		return keventstore.UpdateTwinEvent(event)
	}

	var latestStreetlight model.Streetlight
	err = latestEvent.ToModel(&latestStreetlight)

	if err != nil {
		return err
	}

	if latestStreetlight.PowerState == currentStreetlight.PowerState {
		if currentStreetlight.PowerState == model.PowerOn {
			if isWithDefect(timeNow, latestStreetlight.DateLastSwitchingOn) {
				currentStreetlight.Status = model.LampStatusDefective
			}
			currentStreetlight.DateLastSwitchingOn = timeNow
		}
		if currentStreetlight.PowerState == model.PowerOff {
			if isWithDefect(timeNow, latestStreetlight.DateLastSwitchingOff) {
				currentStreetlight.Status = model.LampStatusDefective
			}
			currentStreetlight.DateLastSwitchingOff = timeNow
		}
	}

	event.SetData(currentStreetlight)
	return keventstore.UpdateTwinEvent(event)
}

// In case of 48h of no change in the state, we consider that lamp with a defect
func isWithDefect(datetimeNow *time.Time, dateLastSwitching *time.Time) bool {
	if dateLastSwitching == nil {
		return false
	}
	timeDifference := datetimeNow.Sub(*dateLastSwitching)
	return timeDifference > time.Hour*48
}
