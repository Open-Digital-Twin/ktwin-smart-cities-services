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

func TestAirQualityObservedServiceSuite(t *testing.T) {

	suite.Run(t, new(AirQualityObservedServiceSuite))
}

type AirQualityObservedServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *AirQualityObservedServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *AirQualityObservedServiceSuite) Test_PoleAirQualityObservedEvent() {
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
				Given new air quality observed event is received
				When all values are under a good AQI level
				Should update event with calculated AQI levels and store it in event store
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-airqualityobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-airqualityobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"CODensity": 8, "NO2Density": 8, "O3Density": 8, "SO2Density": 8, "PM10Density": 8, "PM25Density": 8}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-airqualityobserved-nb001-p00007")
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
					MatchHeader("ce-source", "ngsi-ld-city-airqualityobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-airqualityobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"CODensity":8,"PM10Density":8,"PM25Density":8,"SO2Density":8,"NO2Density":8,"O3Density":8,"COAqiLevel":"MODERATE","PM10AqiLevel":"GOOD","PM25AqiLevel":"GOOD","SO2AqiLevel":"GOOD","O3AqiLevel":"GOOD"}`).
					Reply(http.StatusAccepted)

				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", DEFAULT_UUID).
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "city-pole-nb001-p00007").
					MatchHeader("ce-type", "ktwin.command.city-pole.updateairqualityindex").
					MatchHeader("ce-subject", "").
					BodyString(`{"aqiLevel":"MODERATE"}`).
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
