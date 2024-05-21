package service

import (
	"net/http"
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

func TestParkingSpotServiceSuite(t *testing.T) {

	suite.Run(t, new(ParkingSpotServiceSuite))
}

type ParkingSpotServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *ParkingSpotServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *ParkingSpotServiceSuite) Test_ParkingSpotEvent() {
	defer clock.ResetClockImplementation()
	defer uuid.ResetUuidImplementation()

	uuid.NewUuid = func() string {
		return DEFAULT_UUID
	}

	clock.NowFunc = func() *time.Time {
		now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
		return &now
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
			name: `Empty event`,
			twinEvent: func() *ktwin.TwinEvent {
				return &ktwin.TwinEvent{}
			},
			mockExternalService: func() {},
			expectedError:       nil,
		},
		{
			name: `
				Given new parking spot event is received
				When new parking spot event does not have status
				Should ignore event and log error
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparkingspot"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"status": ""}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-offstreetparkingspot.ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
			},
			expectedError: nil,
		},
		{
			name: `
				Given new parking spot event is received
				When new parking spot has status as occupied
				Should generate command to increment vehicle entrance count
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparkingspot"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"status": "occupied"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-offstreetparkingspot")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", DEFAULT_UUID).
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					MatchHeader("ce-type", "ktwin.command.ngsi-ld-city-offstreetparking.updatevehiclecount").
					MatchHeader("ce-subject", "").
					BodyString(`{"vehicleEntranceCount":1}`).
					Reply(http.StatusAccepted)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new parking spot event is received
				When new parking spot has status as free
				Should generate command to increment vehicle exit count
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparkingspot"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"status": "free"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-offstreetparkingspot")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", DEFAULT_UUID).
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					MatchHeader("ce-type", "ktwin.command.ngsi-ld-city-offstreetparking.updatevehiclecount").
					MatchHeader("ce-subject", "").
					BodyString(`{"vehicleExitCount":1}`).
					Reply(http.StatusAccepted)
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
