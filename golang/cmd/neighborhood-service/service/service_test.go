package service

import (
	"os"
	"testing"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/neighborhood-service/model"
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

func TestNeighborhoodServiceSuite(t *testing.T) {

	suite.Run(t, new(NeighborhoodServiceSuite))
}

type NeighborhoodServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *NeighborhoodServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *NeighborhoodServiceSuite) Test_NeighborhoodEvent() {
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
				Given new command is received and there is no previous event
				When command has aqiLevel GOOD
				Should create neighborhood event with aqiLevel GOOD and dateObserved and dateModified set to current time
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "s4city-city-neighborhood-nb001"
				twinEvent.TwinInterface = "s4city-city-neighborhood"
				twinEvent.CommandName = "updateAirQualityIndex"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"aqiLevel": "GOOD"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("s4city-city-neighborhood-nb001")
				cloudEvent.SetType("ktwin.command.s4city-city-neighborhood.updateAirQualityIndex")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/s4city-city-neighborhood/s4city-city-neighborhood-nb001/latest").
					Reply(404)

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "s4city-city-neighborhood-nb001").
					MatchHeader("ce-type", "ktwin.real.s4city-city-neighborhood").
					MatchHeader("ce-subject", "").
					BodyString(`{"aqiLevel":"GOOD","dateObserved":"2024-01-01T00:00:00Z","dateModified":"2024-01-01T00:00:00Z"}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new command is received and there a previous event
				When the new command has aqiLevel GOOD
				AND previous event has aqiLevel UNHEALTHY
				AND the time difference between the previous event and the current time is less than 60 minutes
				Should update neighborhood event with aqiLevel UNHEALTHY
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "s4city-city-neighborhood-nb001"
				twinEvent.TwinInterface = "s4city-city-neighborhood"
				twinEvent.CommandName = "updateAirQualityIndex"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"aqiLevel": "GOOD"}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("s4city-city-neighborhood-nb001")
				cloudEvent.SetType("ktwin.command.s4city-city-neighborhood.updateAirQualityIndex")
				cloudEvent.SetTime(dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/s4city-city-neighborhood/s4city-city-neighborhood-nb001/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "s4city-city-neighborhood-nb001").
					SetHeader("ce-type", "ktwin.real.s4city-city-neighborhood").
					SetHeader("ce-subject", "").
					JSON(model.Neighborhood{
						AqiLevel:     model.UNHEALTHY,
						DateObserved: dateTime,
					})

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "s4city-city-neighborhood-nb001").
					MatchHeader("ce-type", "ktwin.real.s4city-city-neighborhood").
					MatchHeader("ce-subject", "").
					BodyString(`{"aqiLevel":"UNHEALTHY","dateObserved":"2024-01-01T00:00:00Z","dateModified":"0001-01-01T00:00:00Z"}`).
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
