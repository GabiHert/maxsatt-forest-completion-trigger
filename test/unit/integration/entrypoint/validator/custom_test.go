package validator_test

import (
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/entrypoint/validator"
)

type TestStructNotNil struct {
	Field string `json:"field" validate:"not_nil"`
}

type TestStructOneOfIgnoreCase struct {
	Field string `json:"field" validate:"oneof_ignore_case=apple BANANA Orange"`
}

type TestStructDateTimeFormat struct {
	Field string `json:"field" validate:"datetime_format=2006-01-02T15:04:05Z07:00"`
}

type TestStructHourMinute struct {
	Field string `json:"field" validate:"hour_minute"`
}

func TestCustomValidator_NotNil(t *testing.T) {
	v := validator.CustomValidator()

	test := TestStructNotNil{Field: "value"}
	err := v.Struct("TEST-001", test)
	if err != nil {
		t.Errorf("Struct() with not_nil should not return error, got: %v", err)
	}

	testEmpty := TestStructNotNil{Field: ""}
	err = v.Struct("TEST-001", testEmpty)
	if err != nil {
		t.Errorf("Struct() with not_nil on empty value should not return error, got: %v", err)
	}
}

func TestCustomValidator_OneOfIgnoreCase(t *testing.T) {
	v := validator.CustomValidator()

	tests := []struct {
		name      string
		input     TestStructOneOfIgnoreCase
		wantError bool
	}{
		{
			name:      "exact match lowercase",
			input:     TestStructOneOfIgnoreCase{Field: "apple"},
			wantError: false,
		},
		{
			name:      "exact match uppercase",
			input:     TestStructOneOfIgnoreCase{Field: "BANANA"},
			wantError: false,
		},
		{
			name:      "exact match mixed case",
			input:     TestStructOneOfIgnoreCase{Field: "Orange"},
			wantError: false,
		},
		{
			name:      "case insensitive match - all lowercase",
			input:     TestStructOneOfIgnoreCase{Field: "banana"},
			wantError: false,
		},
		{
			name:      "case insensitive match - all uppercase",
			input:     TestStructOneOfIgnoreCase{Field: "APPLE"},
			wantError: false,
		},
		{
			name:      "case insensitive match - mixed",
			input:     TestStructOneOfIgnoreCase{Field: "oRaNgE"},
			wantError: false,
		},
		{
			name:      "no match",
			input:     TestStructOneOfIgnoreCase{Field: "grape"},
			wantError: true,
		},
		{
			name:      "empty string",
			input:     TestStructOneOfIgnoreCase{Field: ""},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct("TEST-002", tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Struct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCustomValidator_DateTimeFormat(t *testing.T) {
	v := validator.CustomValidator()

	tests := []struct {
		name      string
		input     TestStructDateTimeFormat
		wantError bool
	}{
		{
			name:      "valid datetime",
			input:     TestStructDateTimeFormat{Field: "2025-02-20T10:00:00Z"},
			wantError: false,
		},
		{
			name:      "valid datetime with timezone",
			input:     TestStructDateTimeFormat{Field: "2025-02-20T10:00:00+05:00"},
			wantError: false,
		},
		{
			name:      "empty string is valid",
			input:     TestStructDateTimeFormat{Field: ""},
			wantError: false,
		},
		{
			name:      "invalid format - missing timezone",
			input:     TestStructDateTimeFormat{Field: "2025-02-20T10:00:00"},
			wantError: true,
		},
		{
			name:      "invalid format - wrong format",
			input:     TestStructDateTimeFormat{Field: "2025-02-20 10:00:00"},
			wantError: true,
		},
		{
			name:      "invalid format - date only",
			input:     TestStructDateTimeFormat{Field: "2025-02-20"},
			wantError: true,
		},
		{
			name:      "invalid format - random string",
			input:     TestStructDateTimeFormat{Field: "invalid"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct("TEST-003", tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Struct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCustomValidator_HourMinute(t *testing.T) {
	v := validator.CustomValidator()

	tests := []struct {
		name      string
		input     TestStructHourMinute
		wantError bool
	}{
		{
			name:      "valid hour minute - 00:00",
			input:     TestStructHourMinute{Field: "00:00"},
			wantError: false,
		},
		{
			name:      "valid hour minute - 23:59",
			input:     TestStructHourMinute{Field: "23:59"},
			wantError: false,
		},
		{
			name:      "valid hour minute - 12:30",
			input:     TestStructHourMinute{Field: "12:30"},
			wantError: false,
		},
		{
			name:      "invalid - hour 24",
			input:     TestStructHourMinute{Field: "24:00"},
			wantError: true,
		},
		{
			name:      "invalid - minute 60",
			input:     TestStructHourMinute{Field: "12:60"},
			wantError: true,
		},
		{
			name:      "invalid - negative hour",
			input:     TestStructHourMinute{Field: "-1:30"},
			wantError: true,
		},
		{
			name:      "invalid - negative minute",
			input:     TestStructHourMinute{Field: "12:-1"},
			wantError: true,
		},
		{
			name:      "invalid - no colon",
			input:     TestStructHourMinute{Field: "1230"},
			wantError: true,
		},
		{
			name:      "invalid - multiple colons",
			input:     TestStructHourMinute{Field: "12:30:45"},
			wantError: true,
		},
		{
			name:      "invalid - not numbers",
			input:     TestStructHourMinute{Field: "ab:cd"},
			wantError: true,
		},
		{
			name:      "invalid - empty string",
			input:     TestStructHourMinute{Field: ""},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct("TEST-004", tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Struct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCustomValidator_Struct_ValidInput(t *testing.T) {
	v := validator.CustomValidator()

	type ValidStruct struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"min=0,max=150"`
	}

	validInput := ValidStruct{Name: "John", Age: 30}
	err := v.Struct("TEST-005", validInput)
	if err != nil {
		t.Errorf("Struct() should not return error for valid input, got: %v", err)
	}
}

func TestCustomValidator_Struct_InvalidInput(t *testing.T) {
	v := validator.CustomValidator()

	type InvalidStruct struct {
		Name string `json:"name" validate:"required"`
	}

	invalidInput := InvalidStruct{Name: ""}
	err := v.Struct("TEST-006", invalidInput)
	if err == nil {
		t.Error("Struct() should return error for invalid input")
	}
}

func TestCustomValidator_Struct_GreaterThanValidation(t *testing.T) {
	v := validator.CustomValidator()

	type GreaterThanStruct struct {
		Count int `json:"count" validate:"gt=0"`
	}

	invalidInput := GreaterThanStruct{Count: 0}
	err := v.Struct("TEST-007", invalidInput)
	if err == nil {
		t.Error("Struct() should return error when gt validation fails")
	}

	validInput := GreaterThanStruct{Count: 5}
	err = v.Struct("TEST-007", validInput)
	if err != nil {
		t.Errorf("Struct() should not return error for valid input, got: %v", err)
	}
}

func TestCustomValidator_Struct_NotNilValidation(t *testing.T) {
	v := validator.CustomValidator()

	type NotNilStruct struct {
		Name string `json:"name" validate:"not_nil"`
	}

	input := NotNilStruct{Name: ""}
	err := v.Struct("TEST-008", input)
	if err != nil {
		t.Errorf("Struct() with not_nil should not return error, got: %v", err)
	}
}

func TestCustomValidator_Struct_UnknownTag(t *testing.T) {
	v := validator.CustomValidator()

	type UnknownTagStruct struct {
		Email string `json:"email" validate:"email"`
	}

	invalidInput := UnknownTagStruct{Email: "not-an-email"}
	err := v.Struct("TEST-009", invalidInput)
	if err == nil {
		t.Error("Struct() should return error for invalid email")
	}
}
