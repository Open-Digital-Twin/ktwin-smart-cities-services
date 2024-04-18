package ktwin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/clock"
	"github.com/Open-Digital-Twin/ktwin-smart-cities-services/pkg/logger"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

const (
	EventRealGenerated    = "ktwin.real.%s"
	EventVirtualGenerated = "ktwin.virtual.%s"
	EventCommandExecuted  = "ktwin.command.%s.%s"
)

func GetEventStoreURL() string {
	return os.Getenv("KTWIN_EVENT_STORE")
}

func GetBrokerURL() string {
	fmt.Printf("Broker URL: %s\n", os.Getenv("KTWIN_BROKER"))
	return os.Getenv("KTWIN_BROKER")
}

func PostCloudEvent(event *cloudevents.Event, url string) error {
	if os.Getenv("ENV") == "local" {
		return nil
	}

	client := NewClient()
	response, err := client.Post(url, event)

	if err != nil {
		return errors.New("error to publish cloud event: " + err.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("error to publish cloud event. status code: %d", response.StatusCode)
}

func GetCloudEvent(cloudEvent *cloudevents.Event, url string) (*cloudevents.Event, error) {
	if os.Getenv("ENV") == "local" {
		return cloudEvent, nil
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), GetBrokerURL())

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		logger.NewLogger().Error("failed to create client", err)
		return nil, err
	}

	var event *cloudevents.Event
	if event, err = c.Request(ctx, *cloudEvent); err != nil {
		return nil, errors.New("Error to get Cloud Event: " + err.Error())
	}
	return event, nil

}

type Client struct {
	client http.Client
}

func NewClient() Client {
	return Client{
		client: http.Client{},
	}
}

func (c *Client) Post(url string, event *cloudevents.Event) (*http.Response, error) {
	req, err := c.createRequest(url, event)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *Client) createRequest(url string, cloudEvent *cloudevents.Event) (*http.Request, error) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(cloudEvent.Data()))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ce-id", cloudEvent.ID())
	req.Header.Set("ce-specversion", cloudEvent.SpecVersion())
	req.Header.Set("ce-time", cloudEvent.Time().Format(time.RFC3339))
	req.Header.Set("ce-source", cloudEvent.Source())
	req.Header.Set("ce-type", cloudEvent.Type())
	req.Header.Set("ce-subject", cloudEvent.Subject())

	return req, nil
}

// TwinEvent

type EventType string

const (
	RealEvent    EventType = "real"
	VirtualEvent EventType = "virtual"
	CommandEvent EventType = "command"
)

type TwinEvent struct {
	CloudEvent *cloudevents.Event

	// It is part of CloudEvent type
	EventType     EventType
	TwinInterface string
	CommandName   string

	// The Source of the CloudEvent
	TwinInstance string
}

func NewTwinEvent() *TwinEvent {
	return &TwinEvent{}
}

// Real Event Type: ktwin.real.<twin-interface>
// Virtual Event Type: ktwin.virtual.<twin-interface>
// Command Event Type: ktwin.command.<twin-interface>.<command-name>
func (e *TwinEvent) HandleRequest(r *http.Request) error {
	cloudEvent, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		log.Printf("failed to parse CloudEvent from request: %v", err)
		return err
	}

	ceType := strings.Split(cloudEvent.Type(), ".")
	e.EventType = EventType(ceType[1])
	e.TwinInterface = ceType[2]
	e.TwinInstance = cloudEvent.Source()
	e.CloudEvent = cloudEvent

	if len(ceType) > 3 {
		e.CommandName = ceType[3]
	}

	if e.EventType == "" {
		return errors.New("event type not found")
	}

	if e.EventType == CommandEvent && e.CommandName == "" {
		return errors.New("command name not found")
	}

	return nil
}

func (k *TwinEvent) HandleResponse(r *http.Response) error {
	cloudEvent, err := cloudevents.NewEventFromHTTPResponse(r)
	if err != nil {
		log.Printf("failed to parse CloudEvent from request: %v", err)
		return err
	}
	k.TwinInstance = cloudEvent.Source()
	k.TwinInterface = cloudEvent.Type()
	k.CloudEvent = cloudEvent
	return nil
}

func (k *TwinEvent) ToModel(model interface{}) error {
	err := json.Unmarshal(k.CloudEvent.Data(), model)
	if err != nil {
		return err
	}
	return nil
}

func (k *TwinEvent) SetData(model interface{}) error {
	return k.CloudEvent.SetData(cloudevents.ApplicationJSON, model)
}

func (ktwinEvent *TwinEvent) SetEvent(twinInterface, twinInstance string, data interface{}) {
	ktwinEvent.TwinInterface = twinInterface
	ktwinEvent.TwinInstance = twinInstance
	ceType := fmt.Sprintf("ktwin.%s", ktwinEvent.TwinInterface)
	ktwinEvent.CloudEvent = BuildCloudEvent(ceType, twinInstance, data)
}

func BuildCloudEvent(ceType, ceSource string, data interface{}) *cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetTime(clock.Now())
	event.SetType(ceType)
	event.SetSource(ceSource)
	event.SetData(cloudevents.ApplicationJSON, data)
	return &event
}

type TwinInstanceReference struct {
	Name      string `json:"name"`
	Interface string `json:"interface"`
	Instance  string `json:"instance"`
}

func (t *TwinInstanceReference) ToJSON() string {
	j, _ := json.MarshalIndent(t, "", "    ")
	return string(j)
}

type TwinInstanceGraph struct {
	Name          string                  `json:"name"`
	Interface     string                  `json:"interface"`
	Relationships []TwinInstanceReference `json:"relationships"`
}

func (t *TwinInstanceGraph) ToJSON() string {
	j, _ := json.MarshalIndent(t, "", "    ")
	return string(j)
}

type TwinGraph struct {
	TwinInstancesGraph []TwinInstanceGraph `json:"twinInstances"`
}

func (t *TwinGraph) ToJSON() string {
	j, _ := json.MarshalIndent(t, "", "    ")
	return string(j)
}
