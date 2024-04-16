package ktwin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

func PostCloudEvent(cloudEvent *cloudevents.Event, url string) error {
	ctx := cloudevents.ContextWithTarget(context.Background(), GetBrokerURL())

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		logger.NewLogger().Error("failed to create client", err)
		return err
	}

	if err := c.Send(ctx, *cloudEvent); err != nil {
		return errors.New("Error to publish Cloud Event: " + err.Error())
	}

	return nil
}

func GetCloudEvent(cloudEvent *cloudevents.Event, url string) (*cloudevents.Event, error) {
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
