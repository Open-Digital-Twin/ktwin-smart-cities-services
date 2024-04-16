package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/device-service/model"
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

	// Log event details
	logger.Info(fmt.Sprintf("Event TwinInstance: %s - Event TwinInterface: %s", twinEvent.TwinInstance, twinEvent.TwinInterface))

	err := ktwin.HandleEvent(twinEvent, model.DEVICE_INTERFACE_ID, handleDeviceEvent)

	if err != nil {
		logger.Error("Error processing cloud event request", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing cloud event request"))
		return
	}
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
	logger.Info(fmt.Sprintf("Twin Instance: %s | Twin Interface: %s | %#v", event.TwinInstance, event.TwinInterface, device))

	if device.BatteryLevel != 0 {
		if device.BatteryLevel < BatteryThreshold {
			// Propagate event to real device to measure in low frequency
			device.MeasurementFrequency = LowFrequency
			logger.Info(fmt.Sprintf("Battery Level below threshold. Sending event to real instance: %s", event.TwinInstance))
			err := ktwin.PublishToRealTwin(event.TwinInterface, event.TwinInstance, device)
			if err != nil {
				return err
			}
		} else if device.BatteryLevel > BatteryThreshold {
			// Propagate event to real device to measure in high frequency
			device.MeasurementFrequency = HighFrequency
			logger.Info(fmt.Sprintf("Battery Level above threshold. Sending event to real instance: %s", event.TwinInstance))
			err := ktwin.PublishToRealTwin(event.TwinInterface, event.TwinInstance, device)
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

func main() {
	server.LoadEnv()
	server.StartServer(requestHandler)
}