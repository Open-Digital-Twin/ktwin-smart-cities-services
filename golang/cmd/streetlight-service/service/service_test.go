package service

import (
	"os"
	"testing"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/streetlight-service/model"
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

func TestStreetlightServiceSuite(t *testing.T) {

	suite.Run(t, new(StreetlightServiceSuite))
}

type StreetlightServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *StreetlightServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *StreetlightServiceSuite) Test_StreetlightEvent() {
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
			name: `Empty event`,
			twinEvent: func() *ktwin.TwinEvent {
				return &ktwin.TwinEvent{}
			},
			mockExternalService: func() {},
			expectedError:       nil,
		},
		{
			name: `
				Given event is published and no previous event was published
				When event has powerState "on"
				Should update the event with the new powerState and set the dateLastSwitchingOn
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.RealEvent
				twinEvent.TwinInstance = "ngsi-ld-city-streetlight-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-streetlight"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"powerState": "on"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-streetlight-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-streetlight")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-streetlight/ngsi-ld-city-streetlight-nb001-p00007/latest").
					Reply(404)

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-streetlight-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-streetlight").
					MatchHeader("ce-subject", "").
					BodyString(`{"circuit":"","status":"","powerState":"on","dateLastLampChange":"0001-01-01T00:00:00Z","dateLastSwitchingOn":"2024-01-01T00:00:00Z","dateLastSwitchingOff":"0001-01-01T00:00:00Z","controllingMethod":"","dateServiceStarted":"0001-01-01T00:00:00Z","image":"","annotations":"","lanternHeight":0,"illuminanceLevel":0,"locationCategory":""}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given event is published and no previous event was published
				When event has powerState "off"
				Should update the event with the new powerState and set the dateLastSwitchingOn
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.RealEvent
				twinEvent.TwinInstance = "ngsi-ld-city-streetlight-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-streetlight"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"powerState": "off"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-streetlight-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-streetlight")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-streetlight/ngsi-ld-city-streetlight-nb001-p00007/latest").
					Reply(404)

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-streetlight-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-streetlight").
					MatchHeader("ce-subject", "").
					BodyString(`{"circuit":"","status":"","powerState":"off","dateLastLampChange":"0001-01-01T00:00:00Z","dateLastSwitchingOn":"0001-01-01T00:00:00Z","dateLastSwitchingOff":"2024-01-01T00:00:00Z","controllingMethod":"","dateServiceStarted":"0001-01-01T00:00:00Z","image":"","annotations":"","lanternHeight":0,"illuminanceLevel":0,"locationCategory":""}`).
					Reply(200)

			},
			expectedError: nil,
		},
		{
			name: `
				Given event is published and it has previous event published
				When new event has powerState "off" and latest event has powerState "off" and dateLastSwitchingOff is not older than 2 days
				Should update the event with the new powerState and set the dateLastSwitchingOn
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.RealEvent
				twinEvent.TwinInstance = "ngsi-ld-city-streetlight-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-streetlight"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"powerState": "off"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-streetlight-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-streetlight")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-streetlight/ngsi-ld-city-streetlight-nb001-p00007/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-streetlight-nb001-p00007").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-streetlight").
					SetHeader("ce-subject", "").
					JSON(model.Streetlight{
						PowerState:           "off",
						DateLastSwitchingOff: dateTime,
					})

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-streetlight-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-streetlight").
					MatchHeader("ce-subject", "").
					BodyString(`{"circuit":"","status":"","powerState":"off","dateLastLampChange":"0001-01-01T00:00:00Z","dateLastSwitchingOn":"0001-01-01T00:00:00Z","dateLastSwitchingOff":"2024-01-01T00:00:00Z","controllingMethod":"","dateServiceStarted":"0001-01-01T00:00:00Z","image":"","annotations":"","lanternHeight":0,"illuminanceLevel":0,"locationCategory":""}`).
					Reply(200)

			},
			expectedError: nil,
		},
		{
			name: `
				Given event is published and it has previous event published
				When new event has powerState "off" and latest event has powerState "off" and dateLastSwitchingOff is older than 2 days
				Should update the event with the new powerState and set the dateLastSwitchingOn
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.RealEvent
				twinEvent.TwinInstance = "ngsi-ld-city-streetlight-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-streetlight"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"powerState": "off"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-streetlight-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-streetlight")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				pastDateTime, _ := time.Parse("2006-01-02T15:04:05Z", "2023-01-01T00:00:00Z")
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-streetlight/ngsi-ld-city-streetlight-nb001-p00007/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-streetlight-nb001-p00007").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-streetlight").
					SetHeader("ce-subject", "").
					JSON(model.Streetlight{
						PowerState:           "off",
						DateLastSwitchingOff: pastDateTime,
					})

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-streetlight-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-streetlight").
					MatchHeader("ce-subject", "").
					BodyString(`{"circuit":"","status":"defectiveLamp","powerState":"off","dateLastLampChange":"0001-01-01T00:00:00Z","dateLastSwitchingOn":"0001-01-01T00:00:00Z","dateLastSwitchingOff":"2024-01-01T00:00:00Z","controllingMethod":"","dateServiceStarted":"0001-01-01T00:00:00Z","image":"","annotations":"","lanternHeight":0,"illuminanceLevel":0,"locationCategory":""}`).
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
