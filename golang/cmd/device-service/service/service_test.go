package service

import (
	"os"
	"testing"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/clock"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/ktwin/config"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/uuid"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

const (
	DEFAULT_UUID = "e8e126f6-62fb-40fd-a7cd-8264ca8600d0"
)

func TestDeviceServiceSuite(t *testing.T) {

	suite.Run(t, new(DeviceServiceSuite))
}

type DeviceServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *DeviceServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *DeviceServiceSuite) Test_DeviceEvent() {
	defer clock.ResetClockImplementation()
	defer uuid.ResetUuidImplementation()

	uuid.NewUuid = func() string {
		return DEFAULT_UUID
	}

	clock.NowFunc = func() time.Time {
		now, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-01T00:00:00Z")
		return now
	}
	dateTime := clock.NowFunc()
	dateTimeFormatted := dateTime.Format(time.RFC3339)

	tests := []struct {
		name                string
		mockExternalService func()
		twinEvent           func() *ktwin.TwinEvent
		expectedError       error
	}{
		{
			name: `Invalid Event Type`,
			twinEvent: func() *ktwin.TwinEvent {
				return &ktwin.TwinEvent{}
			},
			mockExternalService: func() {},
			expectedError:       nil,
		},
		{
			name: `
				Given TwinEvent is valid
				When battery level is above threshold
				Should propagate event to real device to measure in low frequency
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.RealEvent
				twinEvent.TwinInstance = "ngsi-ld-city-device-nb001-ofp0003-s0012"
				twinEvent.TwinInterface = "ngsi-ld-city-device"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"batteryLevel": 20}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-device-nb001-ofp0003-s0012")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-device")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-device-nb001-ofp0003-s0012").
					MatchHeader("ce-type", "ktwin.virtual.ngsi-ld-city-device").
					MatchHeader("ce-subject", "").
					BodyString(`{"dataProvider":"","batteryLevel":20,"measurementFrequency":15,"source":"","dateCreated":"0001-01-01T00:00:00Z","dateObserved":"2024-01-01T00:00:00Z","dateModified":"0001-01-01T00:00:00Z"}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given TwinEvent is valid
				When battery level is below threshold
				Should propagate event to real device to measure in high frequency
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.RealEvent
				twinEvent.TwinInstance = "ngsi-ld-city-device-nb001-ofp0003-s0012"
				twinEvent.TwinInterface = "ngsi-ld-city-device"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"batteryLevel": 15}`))
				cloudEvent.SetID("e8e126f6-62fb-40fd-a7cd-8264ca8600d0")
				cloudEvent.SetSource("ngsi-ld-city-device-nb001-ofp0003-s0012")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-device")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-device-nb001-ofp0003-s0012").
					MatchHeader("ce-type", "ktwin.virtual.ngsi-ld-city-device").
					MatchHeader("ce-subject", "").
					BodyString(`{\"dataProvider\":\"\",\"batteryLevel\":20,\"measurementFrequency\":15,\"source\":\"\",\"dateCreated\":\"0001-01-01T00:00:00Z\",\"dateObserved\":\"0001-01-01T00:00:00Z\",\"dateModified\":\"0001-01-01T00:00:00Z\"}`).
					Reply(200)
			},
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			defer gock.Off()
			tt.mockExternalService()

			actualError := HandleEvent(tt.twinEvent())

			s.Assert().Equal(tt.expectedError, actualError)
		})
	}
}
