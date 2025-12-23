package entity

import (
	"testing"
	"time"
)

func TestAnalysis_BasePath(t *testing.T) {
	tests := []struct {
		name               string
		processingConfigID string
		forestID           string
		fieldID            string
		date               time.Time
		expected           string
	}{
		{
			name:               "with processing_config_id",
			processingConfigID: "config-001",
			forestID:           "forest-123",
			fieldID:            "field-456",
			date:               time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected:           "config-001/forest-123/field-456/2024-01-15",
		},
		{
			name:               "without processing_config_id (backward compatibility)",
			processingConfigID: "",
			forestID:           "forest-123",
			fieldID:            "field-456",
			date:               time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected:           "forest-123/field-456/2024-01-15",
		},
		{
			name:               "with UUID processing_config_id",
			processingConfigID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			forestID:           "87428a25-f0d2-4dba-acfd-e3d7320c60002",
			fieldID:            "5a801269-110d-46c0-89dd-09a6dc196002",
			date:               time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected:           "a1b2c3d4-e5f6-7890-abcd-ef1234567890/87428a25-f0d2-4dba-acfd-e3d7320c60002/5a801269-110d-46c0-89dd-09a6dc196002/2024-01-02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analysis{
				ProcessingConfigID: tt.processingConfigID,
				ForestID:           tt.forestID,
				FieldID:            tt.fieldID,
				Date:               tt.date,
			}
			if got := a.BasePath(); got != tt.expected {
				t.Errorf("BasePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAnalysis_DatePath(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "standard date",
			date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "2024-01-15",
		},
		{
			name:     "single digit month and day",
			date:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: "2024-01-02",
		},
		{
			name:     "end of year",
			date:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: "2024-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analysis{Date: tt.date}
			if got := a.DatePath(); got != tt.expected {
				t.Errorf("DatePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}
