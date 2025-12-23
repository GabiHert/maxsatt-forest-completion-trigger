package enums_test

import (
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
)

func TestEventType_Equals(t *testing.T) {
	tests := []struct {
		name      string
		eventType enums.EventType
		input     string
		expected  bool
	}{
		{
			name:      "exact match uppercase",
			eventType: "FOREST_CREATED",
			input:     "FOREST_CREATED",
			expected:  true,
		},
		{
			name:      "case insensitive match lowercase",
			eventType: "FOREST_CREATED",
			input:     "forest_created",
			expected:  true,
		},
		{
			name:      "case insensitive match mixed case",
			eventType: "FOREST_CREATED",
			input:     "Forest_Created",
			expected:  true,
		},
		{
			name:      "no match different value",
			eventType: "FOREST_CREATED",
			input:     "START_ANALYSIS",
			expected:  false,
		},
		{
			name:      "no match empty string",
			eventType: "FOREST_CREATED",
			input:     "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.eventType.Equals(tt.input)
			if result != tt.expected {
				t.Errorf("EventType.Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEventType_String(t *testing.T) {
	tests := []struct {
		name      string
		eventType enums.EventType
		expected  string
	}{
		{
			name:      "uppercase input",
			eventType: "FOREST_CREATED",
			expected:  "FOREST_CREATED",
		},
		{
			name:      "lowercase input converts to uppercase",
			eventType: "forest_created",
			expected:  "FOREST_CREATED",
		},
		{
			name:      "mixed case converts to uppercase",
			eventType: "Forest_Created",
			expected:  "FOREST_CREATED",
		},
		{
			name:      "empty string",
			eventType: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.eventType.String()
			if result != tt.expected {
				t.Errorf("EventType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}
