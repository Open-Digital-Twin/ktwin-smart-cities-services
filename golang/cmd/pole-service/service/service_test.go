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

func TestPoleServiceSuite(t *testing.T) {

	suite.Run(t, new(PoleServiceSuite))
}

type PoleServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *PoleServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *PoleServiceSuite) Test_PoleAirQualityObservedEvent() {
	defer clock.ResetClockImplementation()
	defer uuid.ResetUuidImplementation()

	uuid.NewUuid = func() string {
		return DEFAULT_UUID
	}

	clock.NowFunc = func() *time.Time {
		now, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-01T00:00:00Z")
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
				When new parking spot has status as occupied
				Should generate command to increment vehicle entrance count
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-airqualityobserved-nb001-p00037"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparkingspot"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"status": "occupied"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-airqualityobserved-nb001-p00037")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-airqualityobserved")
				cloudEvent.SetTime(*dateTime)

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
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparkingspot-nb001-ofp0005-s0008").
					MatchHeader("ce-type", "ktwin.command.ngsi-ld-city-offstreetparking.updatevehiclecount").
					MatchHeader("ce-subject", "").
					BodyString(`{"vehicleEntranceCount":1}`).
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

func (s *PoleServiceSuite) Test_PoleCrowdFlowObservedEvent() {
	defer clock.ResetClockImplementation()
	defer uuid.ResetUuidImplementation()

	uuid.NewUuid = func() string {
		return DEFAULT_UUID
	}

	clock.NowFunc = func() *time.Time {
		now, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-01T00:00:00Z")
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
				Given new crowd flow observed event is received
				When average crowd speed is below threshold AND average headway time is below threshold
				Should update event as congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-crowdflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-crowdflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageCrowdSpeed": 2, "averageHeadwayTime": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-crowdflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-crowdflowobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-crowdflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-crowdflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageCrowdSpeed":2,"congested":true,"averageHeadwayTime":1}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new crowd flow observed event is received
				When average crowd speed is below threshold AND average headway time is above threshold
				Should update event as congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-crowdflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-crowdflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageCrowdSpeed": 2, "averageHeadwayTime": 2}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-crowdflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-crowdflowobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-crowdflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-crowdflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageCrowdSpeed":2,"congested":true,"averageHeadwayTime":2}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new crowd flow observed event is received
				When average crowd speed is above threshold AND average headway time is below threshold
				Should update event as congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-crowdflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-crowdflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageCrowdSpeed": 4, "averageHeadwayTime": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-crowdflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-crowdflowobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-crowdflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-crowdflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageCrowdSpeed":4,"congested":true,"averageHeadwayTime":1}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new crowd flow observed event is received
				When average crowd speed is above threshold AND average headway time is above threshold
				Should update event as not congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-crowdflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-crowdflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageCrowdSpeed": 5, "averageHeadwayTime": 3}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-crowdflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-crowdflowobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-crowdflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-crowdflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageCrowdSpeed":5,"congested":false,"averageHeadwayTime":3}`).
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
