package service

// func strPtr(s string) *string {
// 	return &s
// }

// func TestDeviceServiceSuite(t *testing.T) {

// 	suite.Run(t, new(DeviceServiceSuite))
// }

// type DeviceServiceSuite struct {
// 	suite.Suite

// 	brokerUrl     string
// 	eventStoreUrl string
// }

// func (s *DeviceServiceSuite) SetupSuite() {
// 	os.Setenv("ENV", "test")
// 	config.LoadEnv()

// 	s.brokerUrl = os.Getenv("KTWIN_BROKER")
// 	s.eventStoreUrl = os.Getenv("KTWIN_EVENT_STORE")
// }

// func (s *DeviceServiceSuite) Test_DeviceEvent() {
// 	dateTime, _ := time.Parse("2006-01-02T15:04:05Z", "2024-01-01T00:00:00Z")

// 	tests := []struct {
// 		name                string
// 		mockExternalService func()
// 		twinEvent           *ktwin.TwinEvent
// 		expectedError       error
// 	}{
// 		{
// 			name:                `Invalid Event Type`,
// 			twinEvent:           &ktwin.TwinEvent{},
// 			mockExternalService: func() {},
// 			expectedError:       nil,
// 		},
// 		{
// 			name: `Invalid Event Type`,
// 			twinEvent: &ktwin.TwinEvent{
// 				EventType:     ktwin.RealEvent,
// 				TwinInstance:  "ngsi-ld-city-device",
// 				TwinInterface: "ngsi-ld-city-device",
// 				CloudEvent: &cloudevents.Event{
// 					Context: &cloudevents.EventContextV1{
// 						ID:              "event-id",
// 						DataContentType: strPtr("application/json"),
// 						Source:          *cloudevents.ParseURIRef("ngsi-ld-city-device"),
// 						Type:            "ktwin.real.ngsi-ld-city-device",
// 						Time:            &cloudevents.Timestamp{Time: dateTime},
// 						DataSchema:      cloudevents.ParseURI(""),
// 					},
// 					DataEncoded: []byte(`{"batteryLevel": 20}`),
// 				},
// 			},
// 			mockExternalService: func() {
// 				gock.New(s.brokerUrl).Post("/").BodyString(`{"batteryLevel": 20}`).Reply(200)
// 				//gock.New(s.CONFIG_URL).
// 				// queryConfigByConfigId, _ := json.Marshal(repository.GetConfigDetailsByConfigIdQuery(123))
// 				// gock.New(s.CONFIG_URL).
// 				// 	BodyString(string(queryConfigByConfigId)).
// 				// 	Reply(http.StatusOK).
// 				// 	JSON(repository.ExternalConfig{
// 				// 		Data: repository.Data{
// 				// 			TravelConfigById: datafaker.CreateFullExternalConfig(dateTime),
// 				// 		},
// 				// 	})

// 				// gock.New(s.AUTH_PERMISSIONS_URL).
// 				// 	BodyString(`{"travelConfigUuid": "Uuid"}`).
// 				// 	Reply(http.StatusOK).
// 				// 	JSON(authRepository.DataAccessResponse{})
// 			},
// 			expectedError: nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		s.Run(tt.name, func() {
// 			defer gock.Off()
// 			httptest.NewServer(nil)
// 			tt.mockExternalService()

// 			actualError := HandleEvent(tt.twinEvent)

// 			s.Assert().Equal(tt.expectedError, actualError)
// 		})
// 	}
// }
