package service

import (
	"fmt"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/device-service/model"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/kevent"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/keventstore"
	log "github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
)

var logger = log.NewLogger()

func HandleEvent(event *ktwin.TwinEvent) error {
	return kevent.HandleEvent(event, model.TWIN_INTERFACE_DEVICE, handleDeviceEvent)
}

func handleDeviceEvent(event *ktwin.TwinEvent) error {
	const (
		HighFrequency    = 15 // 15 min
		LowFrequency     = 60 // 60 min
		BatteryThreshold = 15 // percentage of battery available
	)

	device := model.Device{}
	err := event.ToModel(&device)

	if err != nil {
		logger.Error("Error parsing event data", err)
		return nil
	}

	device.DateObserved = time.Now()
	logger.Info(fmt.Sprintf("CloudEvent: %v", string(event.CloudEvent.DataEncoded)))

	if device.BatteryLevel != 0 {
		if device.BatteryLevel < BatteryThreshold {
			// Propagate event to real device to measure in low frequency
			device.MeasurementFrequency = LowFrequency
			logger.Info(fmt.Sprintf("Battery Level below threshold. Sending event to real instance: %s", event.TwinInstance))
			err := kevent.PublishToRealTwin(event.TwinInterface, event.TwinInstance, device)
			if err != nil {
				return err
			}
		} else if device.BatteryLevel > BatteryThreshold {
			// Propagate event to real device to measure in high frequency
			device.MeasurementFrequency = HighFrequency
			logger.Info(fmt.Sprintf("Battery Level above threshold. Sending event to real instance: %s", event.TwinInstance))
			err := kevent.PublishToRealTwin(event.TwinInterface, event.TwinInstance, device)
			if err != nil {
				return err
			}
		}
		event.SetData(device)
		keventstore.UpdateTwinEvent(event)
	} else {
		logger.Info("Battery level was not provided")
	}

	return nil
}
