package service

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/weather-observed-service/model"
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

func TestWeatherObservedServiceSuite(t *testing.T) {

	suite.Run(t, new(WeatherObservedServiceSuite))
}

type WeatherObservedServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *WeatherObservedServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *WeatherObservedServiceSuite) Test_WeatherObservedEvent() {
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
				Given new weather observed event is received
				When there is not latest event
				Should consider current event as latest and store it in event store
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-weatherobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-weatherobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"atmosphericPressure": 10, "temperature": 8, "relativeHumidity": 8, "windSpeed": 8}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-weatherobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-weatherobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-weatherobserved/ngsi-ld-city-weatherobserved-nb001-p00007/latest").
					Reply(http.StatusNotFound)

				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-weatherobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"pressureTendency":"steady","atmosphericPressure":10,"dewpoint":-10.399999999999999,"feelsLikeTemperature":-1.9253082357521691,"temperature":8,"relativeHumidity":8,"windSpeed":8}`).
					Reply(http.StatusOK)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new weather observed event is received
				When there is a latest event in event store AND athmospheric pressure is higher than previous event
				Should update the event with calculated values and store it in event store
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-weatherobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-weatherobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"atmosphericPressure": 10, "temperature": 8, "relativeHumidity": 8, "windSpeed": 8}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-weatherobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-weatherobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-weatherobserved/ngsi-ld-city-weatherobserved-nb001-p00007/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
					SetHeader("ce-subject", "").
					JSON(model.WeatherObservedEvent{
						AtmosphericPressure: 2,
						Temperature:         2,
						RelativeHumidity:    2,
						WindSpeed:           2,
					})

				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-weatherobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"pressureTendency":"raising","atmosphericPressure":10,"dewpoint":-10.399999999999999,"feelsLikeTemperature":-1.9253082357521691,"temperature":8,"relativeHumidity":8,"windSpeed":8}`).
					Reply(http.StatusOK)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new weather observed event is received
				When there is a latest event in event store AND atmospheric pressure is lower than previous event
				Should update the event with calculated values and store it in event store
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-weatherobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-weatherobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"atmosphericPressure": 10, "temperature": 8, "relativeHumidity": 8, "windSpeed": 8}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-weatherobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-weatherobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-weatherobserved/ngsi-ld-city-weatherobserved-nb001-p00007/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
					SetHeader("ce-subject", "").
					JSON(model.WeatherObservedEvent{
						AtmosphericPressure: 12,
						Temperature:         2,
						RelativeHumidity:    2,
						WindSpeed:           2,
					})

				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-weatherobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"pressureTendency":"falling","atmosphericPressure":10,"dewpoint":-10.399999999999999,"feelsLikeTemperature":-1.9253082357521691,"temperature":8,"relativeHumidity":8,"windSpeed":8}`).
					Reply(http.StatusOK)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new weather observed event is received
				When there is a latest event in event store AND atmospheric pressure is equal of previous event
				Should update the event with calculated values and store it in event store
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-weatherobserved-nb001-p00007"
				twinEvent.TwinInterface = "ngsi-ld-city-weatherobserved"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"atmosphericPressure": 10, "temperature": 8, "relativeHumidity": 8, "windSpeed": 8}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-weatherobserved-nb001-p00007")
				cloudEvent.SetType("ktwin.real.ngsi-ld-city-weatherobserved")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-weatherobserved/ngsi-ld-city-weatherobserved-nb001-p00007/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
					SetHeader("ce-subject", "").
					JSON(model.WeatherObservedEvent{
						AtmosphericPressure: 10,
						Temperature:         2,
						RelativeHumidity:    2,
						WindSpeed:           2,
					})

				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.store.ngsi-ld-city-weatherobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"pressureTendency":"steady","atmosphericPressure":10,"dewpoint":-10.399999999999999,"feelsLikeTemperature":-1.9253082357521691,"temperature":8,"relativeHumidity":8,"windSpeed":8}`).
					Reply(http.StatusOK)
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
