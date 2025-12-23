package dto

import (
	"time"
)

type Event struct {
	EventDate          time.Time   `json:"event_date"`
	EventData          *ForestData `json:"event_data"`
	EventId            string      `json:"event_id"`
	EventCorrelationId string      `json:"event_correlation_id"`
	BusinessUniqueId   string      `json:"business_unique_id"`
	Source             string      `json:"source"`
	EventType          string      `json:"event_type" validate:"required"`
	SpecVersion        int         `json:"spec_version"`
}
