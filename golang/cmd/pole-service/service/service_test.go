package service

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/pole-service/model"
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
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-airqualityobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-airqualityobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"CODensity":8,"PM10Density":8,"PM25Density":8,"SO2Density":8,"NO2Density":8,"O3Density":8,"COAqiLevel":"MODERATE","PM10AqiLevel":"GOOD","PM25AqiLevel":"GOOD","SO2AqiLevel":"GOOD","O3AqiLevel":"GOOD"}`).
					Reply(200)

				gock.New(s.brokerUrl).
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "city-pole-nb001-p00007").
					MatchHeader("ce-type", "ktwin.command.s4city-city-neighborhood.updateairqualityindex").
					MatchHeader("ce-subject", "").
					BodyString(`{"aqiLevel":"MODERATE"}`).
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

func (s *PoleServiceSuite) Test_PoleTrafficFlowObservedEvent() {
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
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":3,"congested":true,"averageHeadwayTime":1}`).
					Reply(200)
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
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":3,"congested":true,"averageHeadwayTime":3}`).
					Reply(200)
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
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":13,"congested":true,"averageHeadwayTime":1}`).
					Reply(200)
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
				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-trafficflowobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-trafficflowobserved").
					MatchHeader("ce-subject", "").
					BodyString(`{"averageVehicleSpeed":13,"congested":false,"averageHeadwayTime":3}`).
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

func (s *PoleServiceSuite) Test_PoleWeatherObservedEvent() {
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

				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
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

				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
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

				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
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

				gock.New(s.eventStoreUrl).
					Post("/api/v1/twin-events").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-weatherobserved-nb001-p00007").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-weatherobserved").
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
