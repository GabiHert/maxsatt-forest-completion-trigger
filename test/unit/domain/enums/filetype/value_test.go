package filetype_test

import (
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums/filetype"
)

func TestGetFileType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		isValid  bool
	}{
		{
			name:     "valid CLIMATE uppercase",
			input:    "CLIMATE",
			expected: "climate",
			isValid:  true,
		},
		{
			name:     "valid climate lowercase",
			input:    "climate",
			expected: "climate",
			isValid:  true,
		},
		{
			name:     "valid Climate mixed case",
			input:    "Climate",
			expected: "climate",
			isValid:  true,
		},
		{
			name:     "valid RAW uppercase",
			input:    "RAW",
			expected: "raw",
			isValid:  true,
		},
		{
			name:     "valid raw lowercase",
			input:    "raw",
			expected: "raw",
			isValid:  true,
		},
		{
			name:     "valid CLEAN uppercase",
			input:    "CLEAN",
			expected: "clean",
			isValid:  true,
		},
		{
			name:     "valid clean lowercase",
			input:    "clean",
			expected: "clean",
			isValid:  true,
		},
		{
			name:     "valid TIF uppercase",
			input:    "TIF",
			expected: "tif",
			isValid:  true,
		},
		{
			name:     "valid tif lowercase",
			input:    "tif",
			expected: "tif",
			isValid:  true,
		},
		{
			name:     "invalid file type defaults to unknown",
			input:    "INVALID_TYPE",
			expected: "unknown",
			isValid:  false,
		},
		{
			name:     "empty string defaults to unknown",
			input:    "",
			expected: "unknown",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filetype.GetFileType(tt.input)

			resultStr := string(result)
			if resultStr != tt.expected {
				t.Errorf("GetFileType() = %v, want %v", resultStr, tt.expected)
			}

			if tt.isValid {
				switch tt.input {
				case "CLIMATE", "climate", "Climate":
					if result != filetype.Climate {
						t.Errorf("GetFileType() should return Climate constant")
					}
				case "RAW", "raw":
					if result != filetype.Raw {
						t.Errorf("GetFileType() should return Raw constant")
					}
				case "CLEAN", "clean":
					if result != filetype.Clean {
						t.Errorf("GetFileType() should return Clean constant")
					}
				case "TIF", "tif":
					if result != filetype.Tif {
						t.Errorf("GetFileType() should return Tif constant")
					}
				}
			} else {
				if result != filetype.Unknown {
					t.Errorf("GetFileType() should return Unknown constant for invalid input")
				}
			}
		})
	}
}

func TestFileTypeConstants(t *testing.T) {
	if filetype.Climate.String() != "CLIMATE" {
		t.Errorf("Climate = %v, want CLIMATE", filetype.Climate)
	}

	if filetype.Raw.String() != "RAW" {
		t.Errorf("Raw = %v, want RAW", filetype.Raw)
	}

	if filetype.Clean.String() != "CLEAN" {
		t.Errorf("Clean = %v, want CLEAN", filetype.Clean)
	}

	if filetype.Tif.String() != "TIF" {
		t.Errorf("Tif = %v, want TIF", filetype.Tif)
	}

	if filetype.Unknown.String() != "UNKNOWN" {
		t.Errorf("Unknown = %v, want UNKNOWN", filetype.Unknown)
	}
}
