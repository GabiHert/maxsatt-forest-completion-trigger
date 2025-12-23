package eventtype_test

import (
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/eventtype"
)

func TestGetEventType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		isValid  bool
	}{
		{
			name:     "valid FOREST_CREATED uppercase",
			input:    "FOREST_CREATED",
			expected: "forest_created",
			isValid:  true,
		},
		{
			name:     "valid forest_created lowercase",
			input:    "forest_created",
			expected: "forest_created",
			isValid:  true,
		},
		{
			name:     "valid Forest_Created mixed case",
			input:    "Forest_Created",
			expected: "forest_created",
			isValid:  true,
		},
		{
			name:     "valid START_ANALYSIS uppercase",
			input:    "START_ANALYSIS",
			expected: "start_analysis",
			isValid:  true,
		},
		{
			name:     "valid start_analysis lowercase",
			input:    "start_analysis",
			expected: "start_analysis",
			isValid:  true,
		},
		{
			name:     "valid READY_ANALYSIS uppercase",
			input:    "READY_ANALYSIS",
			expected: "ready_analysis",
			isValid:  true,
		},
		{
			name:     "valid ready_analysis lowercase",
			input:    "ready_analysis",
			expected: "ready_analysis",
			isValid:  true,
		},
		{
			name:     "valid Ready_Analysis mixed case",
			input:    "Ready_Analysis",
			expected: "ready_analysis",
			isValid:  true,
		},
		{
			name:     "invalid event type",
			input:    "INVALID_TYPE",
			expected: "unknown",
			isValid:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "unknown",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := eventtype.GetEventType(tt.input)

			resultStr := string(result)
			if resultStr != tt.expected {
				t.Errorf("GetEventType() = %v, want %v", resultStr, tt.expected)
			}

			if tt.isValid {
				switch tt.input {
				case "FOREST_CREATED", "forest_created", "Forest_Created":
					if result != eventtype.ForestCreated {
						t.Errorf("GetEventType() should return ForestCreated constant")
					}
				case "START_ANALYSIS", "start_analysis":
					if result != eventtype.StartAnalysis {
						t.Errorf("GetEventType() should return StartAnalysis constant")
					}
				case "READY_ANALYSIS", "ready_analysis", "Ready_Analysis":
					if result != eventtype.ReadyAnalysis {
						t.Errorf("GetEventType() should return ReadyAnalysis constant")
					}
				}
			} else {
				if result != eventtype.Unknown {
					t.Errorf("GetEventType() should return Unknown constant for invalid input")
				}
			}
		})
	}
}

func TestEventTypeConstants(t *testing.T) {
	if eventtype.ForestCreated.String() != "FOREST_CREATED" {
		t.Errorf("ForestCreated = %v, want FOREST_CREATED", eventtype.ForestCreated)
	}

	if eventtype.StartAnalysis.String() != "START_ANALYSIS" {
		t.Errorf("StartAnalysis = %v, want START_ANALYSIS", eventtype.StartAnalysis)
	}

	if eventtype.ReadyAnalysis.String() != "READY_ANALYSIS" {
		t.Errorf("ReadyAnalysis = %v, want READY_ANALYSIS", eventtype.ReadyAnalysis)
	}

	if eventtype.Unknown.String() != "UNKNOWN" {
		t.Errorf("Unknown = %v, want UNKNOWN", eventtype.Unknown)
	}
}
