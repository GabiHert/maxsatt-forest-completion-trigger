package enums_test

import (
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
)

func TestFileType_Equals(t *testing.T) {
	tests := []struct {
		name     string
		fileType enums.FileType
		input    string
		expected bool
	}{
		{
			name:     "exact match uppercase",
			fileType: "CLIMATE",
			input:    "CLIMATE",
			expected: true,
		},
		{
			name:     "case insensitive match lowercase",
			fileType: "CLIMATE",
			input:    "climate",
			expected: true,
		},
		{
			name:     "case insensitive match mixed case",
			fileType: "CLIMATE",
			input:    "Climate",
			expected: true,
		},
		{
			name:     "no match different value",
			fileType: "CLIMATE",
			input:    "RAW",
			expected: false,
		},
		{
			name:     "no match empty string",
			fileType: "CLIMATE",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fileType.Equals(tt.input)
			if result != tt.expected {
				t.Errorf("FileType.Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileType_String(t *testing.T) {
	tests := []struct {
		name     string
		fileType enums.FileType
		expected string
	}{
		{
			name:     "uppercase input",
			fileType: "CLIMATE",
			expected: "CLIMATE",
		},
		{
			name:     "lowercase input converts to uppercase",
			fileType: "climate",
			expected: "CLIMATE",
		},
		{
			name:     "mixed case converts to uppercase",
			fileType: "Climate",
			expected: "CLIMATE",
		},
		{
			name:     "empty string",
			fileType: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fileType.String()
			if result != tt.expected {
				t.Errorf("FileType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}
