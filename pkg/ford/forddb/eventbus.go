package forddb

import (
	"reflect"
)

type EventID struct{ StringResourceID[*Event] }

type EventTypeID string

type EventType[T EventPayload] struct {
	ID          EventTypeID
	PayloadType reflect.Type
}

type EventPayload interface {
}

type RawEventPayload = map[string]interface{}

type Event struct {
	ResourceMetadata[EventID, *Event] `json:"metadata"`

	Type    EventTypeID `json:"event_type"`
	Payload any         `json:"event_payload"`
}

type EventBus struct {
}
