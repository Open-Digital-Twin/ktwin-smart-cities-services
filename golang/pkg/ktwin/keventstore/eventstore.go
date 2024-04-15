package keventstore

import (
	"fmt"
	"net/http"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
)

func GetLatestTwinEvent(twinInterface, twinInstance string) (*ktwin.TwinEvent, error) {
	url := fmt.Sprintf("%s/api/v1/twin-events/%s/%s/latest", ktwin.GetEventStoreURL(), twinInterface, twinInstance)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	event := ktwin.NewTwinEvent()
	err = event.HandleResponse(response)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func UpdateTwinEvent(twinEvent *ktwin.TwinEvent) error {
	url := ktwin.GetEventStoreURL() + "/api/v1/twin-events"
	return ktwin.PostCloudEvent(twinEvent.CloudEvent, url)
}
