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

func TestTrafficFlowObservedServiceSuite(t *testing.T) {

	suite.Run(t, new(TrafficFlowObservedServiceSuite))
}

type TrafficFlowObservedServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *TrafficFlowObservedServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *TrafficFlowObservedServiceSuite) Test_TrafficFlowObservedEvent() {
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
				Given new traffic flow observed event is received
				When average vehicle speed is below threshold AND average headway time is below threshold
				Should update event as congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-trafficflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-trafficflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageVehicleSpeed": 3, "averageHeadwayTime": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-trafficflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-trafficflowobserved")
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
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":3,"congested":true,"averageHeadwayTime":1}`).
					Reply(http.StatusAccepted)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new traffic flow observed event is received
				When average vehicle speed is below threshold AND average headway time is above threshold
				Should update event as congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-trafficflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-trafficflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageVehicleSpeed": 3, "averageHeadwayTime": 3}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-trafficflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-trafficflowobserved")
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
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":3,"congested":true,"averageHeadwayTime":3}`).
					Reply(http.StatusAccepted)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new traffic flow observed event is received
				When average vehicle speed is above threshold AND average headway time is below threshold
				Should update event as congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-trafficflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-trafficflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageVehicleSpeed": 13, "averageHeadwayTime": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-trafficflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-trafficflowobserved")
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
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":13,"congested":true,"averageHeadwayTime":1}`).
					Reply(http.StatusAccepted)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new traffic flow observed event is received
				When average vehicle speed is above threshold AND average headway time is above threshold
				Should update event as not congested
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-trafficflowobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-trafficflowobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"averageVehicleSpeed": 13, "averageHeadwayTime": 3}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-trafficflowobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-trafficflowobserved")
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
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":13,"congested":false,"averageHeadwayTime":3}`).
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
