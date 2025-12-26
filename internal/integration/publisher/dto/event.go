package dto

import (
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"

	"github.com/google/uuid"
)

type Event struct {
	EventDate          time.Time `json:"event_date"`
	EventData          any       `json:"event_data"`
	EventId            string    `json:"event_id"`
	EventCorrelationId string    `json:"event_correlation_id"`
	BusinessUniqueId   string    `json:"business_unique_id,omitempty"`
	Source             string    `json:"source"`
	EventType          string    `json:"event_type"`
	SpecVersion        int       `json:"spec_version"`
}

// NewEvent creates a new Event with the default service name as source.
func NewEvent(correlationId string, eventType enums.EventType, eventData any) Event {
	return Event{
		SpecVersion:        1,
		EventId:            uuid.NewString(),
		EventCorrelationId: correlationId,
		EventDate:          time.Now().UTC(),
		Source:             properties.Properties().Application.ServiceName,
		EventType:          eventType.String(),
		EventData:          eventData,
	}
}

// NewEventWithSource creates a new Event with a custom source.
func NewEventWithSource(correlationId string, source string, eventType enums.EventType, eventData any) Event {
	return Event{
		SpecVersion:        1,
		EventId:            uuid.NewString(),
		EventCorrelationId: correlationId,
		EventDate:          time.Now().UTC(),
		Source:             source,
		EventType:          eventType.String(),
		EventData:          eventData,
	}
}
