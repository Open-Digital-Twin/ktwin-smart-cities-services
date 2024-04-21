package service

import (
	"os"
	"testing"
	"time"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/cmd/parking-service/model"
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

func TestParkingServiceSuite(t *testing.T) {

	suite.Run(t, new(ParkingServiceSuite))
}

type ParkingServiceSuite struct {
	suite.Suite

	brokerUrl     string
	eventStoreUrl string
}

func (s *ParkingServiceSuite) SetupSuite() {
	os.Setenv("ENV", "test")
	config.LoadEnv()

	s.brokerUrl = os.Getenv("KTWIN_BROKER")
	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
}

func (s *ParkingServiceSuite) Test_ParkingEvent() {
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
				Given new command is received and there is no previous event
				When command vehicleEntranceCount and vehicleExitCount are 0
				Should not create parking event
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparking-nb001-ofp0005"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparking"
				twinEvent.CommandName = "updateVehicleCount"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"vehicleEntranceCount": 0}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparking-nb001-ofp0005")
				cloudEvent.SetType("ktwin.command.ngsi-ld-city-offstreetparking.updateVehicleCount")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {},
			expectedError:       nil,
		},
		{
			name: `
				Given new command is received and there is no previous event
				When command has a valid vehicleEntranceCount property
				Should create parking event and set OccupiedSpotNumber and TotalSpotNumber fields
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparking-nb001-ofp0005"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparking"
				twinEvent.CommandName = "updateVehicleCount"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"vehicleEntranceCount": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparking-nb001-ofp0005")
				cloudEvent.SetType("ktwin.command.ngsi-ld-city-offstreetparking.updateVehicleCount")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-offstreetparking/ngsi-ld-city-offstreetparking-nb001-ofp0005/latest").
					Reply(404)

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-offstreetparking").
					MatchHeader("ce-subject", "").
					BodyString(`{"occupiedSpotNumber":1,"totalSpotNumber":50}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new command is received and there is no previous event
				When command has a valid vehicleExitCount property
				Should create parking event and set OccupiedSpotNumber and TotalSpotNumber fields
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparking-nb001-ofp0005"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparking"
				twinEvent.CommandName = "updateVehicleCount"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"vehicleExitCount": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparking-nb001-ofp0005")
				cloudEvent.SetType("ktwin.command.ngsi-ld-city-offstreetparking.updateVehicleCount")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-offstreetparking/ngsi-ld-city-offstreetparking-nb001-ofp0005/latest").
					Reply(404)

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-offstreetparking").
					MatchHeader("ce-subject", "").
					BodyString(`{"occupiedSpotNumber":0,"totalSpotNumber":50}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
				Given new command is received and there is previous event
				When command has a valid vehicleExitCount property
				Should update parking event and decrement OccupiedSpotNumber field
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparking-nb001-ofp0005"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparking"
				twinEvent.CommandName = "updateVehicleCount"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"vehicleExitCount": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparking-nb001-ofp0005")
				cloudEvent.SetType("ktwin.command.ngsi-ld-city-offstreetparking.updateVehicleCount")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-offstreetparking/ngsi-ld-city-offstreetparking-nb001-ofp0005/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-offstreetparking").
					SetHeader("ce-subject", "").
					JSON(model.OffStreetParking{
						OccupiedSpotNumber: 1,
						TotalSpotNumber:    50,
					})

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-offstreetparking").
					MatchHeader("ce-subject", "").
					BodyString(`{"occupiedSpotNumber":0,"totalSpotNumber":50}`).
					Reply(200)
			},
			expectedError: nil,
		},
		{
			name: `
			Given new command is received and there is previous event
			When command has a valid vehicleEntranceCount property
			Should update parking event and increment OccupiedSpotNumber field
			`,
			twinEvent: func() *ktwin.TwinEvent {
				twinEvent := ktwin.NewTwinEvent()
				twinEvent.EventType = ktwin.CommandEvent
				twinEvent.TwinInstance = "ngsi-ld-city-offstreetparking-nb001-ofp0005"
				twinEvent.TwinInterface = "ngsi-ld-city-offstreetparking"
				twinEvent.CommandName = "updateVehicleCount"

				cloudEvent := cloudevents.NewEvent()
				cloudEvent.SetData("application/json", []byte(`{"vehicleEntranceCount": 1}`))
				cloudEvent.SetID("")
				cloudEvent.SetSource("ngsi-ld-city-offstreetparking-nb001-ofp0005")
				cloudEvent.SetType("ktwin.command.ngsi-ld-city-offstreetparking.updateVehicleCount")
				cloudEvent.SetTime(*dateTime)

				twinEvent.CloudEvent = &cloudEvent
				return twinEvent
			},
			mockExternalService: func() {
				gock.New(s.eventStoreUrl).
					Get("/api/v1/twin-events/ngsi-ld-city-offstreetparking/ngsi-ld-city-offstreetparking-nb001-ofp0005/latest").
					Reply(200).
					SetHeader("Content-Type", "application/json").
					SetHeader("ce-specversion", "1.0").
					SetHeader("ce-time", dateTimeFormatted).
					SetHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					SetHeader("ce-type", "ktwin.real.ngsi-ld-city-offstreetparking").
					SetHeader("ce-subject", "").
					JSON(model.OffStreetParking{
						OccupiedSpotNumber: 1,
						TotalSpotNumber:    50,
					})

				gock.New(s.eventStoreUrl+"/api/v1/twin-events").
					Post("/").
					MatchHeader("Content-Type", "application/json").
					MatchHeader("ce-id", "").
					MatchHeader("ce-specversion", "1.0").
					MatchHeader("ce-time", dateTimeFormatted).
					MatchHeader("ce-source", "ngsi-ld-city-offstreetparking-nb001-ofp0005").
					MatchHeader("ce-type", "ktwin.real.ngsi-ld-city-offstreetparking").
					MatchHeader("ce-subject", "").
					BodyString(`{"occupiedSpotNumber":2,"totalSpotNumber":50}`).
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
