package error_test

import (
	"strings"
	"testing"

	domainError "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/error"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

func TestNewEmptyDatasetError(t *testing.T) {
	err := domainError.NewEmptyDatasetError()

	if err == nil {
		t.Error("NewEmptyDatasetError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewEmptyDatasetError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "delta dataset is empty") {
		t.Errorf("NewEmptyDatasetError() description should contain 'delta dataset is empty', got: %s", description)
	}
}

func TestNewMalformedDataError(t *testing.T) {
	err := domainError.NewMalformedDataError()

	if err == nil {
		t.Error("NewMalformedDataError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewMalformedDataError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "malformed data") {
		t.Errorf("NewMalformedDataError() description should contain 'malformed data', got: %s", description)
	}
}

func TestNewNilMetadataError(t *testing.T) {
	err := domainError.NewNilMetadataError()

	if err == nil {
		t.Error("NewNilMetadataError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewNilMetadataError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "metadata cannot be nil") {
		t.Errorf("NewNilMetadataError() description should contain 'metadata cannot be nil', got: %s", description)
	}
}

func TestNewNilWeatherMetricsError(t *testing.T) {
	err := domainError.NewNilWeatherMetricsError()

	if err == nil {
		t.Error("NewNilWeatherMetricsError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewNilWeatherMetricsError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "weather metrics cannot be nil") {
		t.Errorf("NewNilWeatherMetricsError() description should contain 'weather metrics cannot be nil', got: %s", description)
	}
}

func TestNewNoDatesFoundError(t *testing.T) {
	err := domainError.NewNoDatesFoundError()

	if err == nil {
		t.Error("NewNoDatesFoundError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewNoDatesFoundError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "no dates found") {
		t.Errorf("NewNoDatesFoundError() description should contain 'no dates found', got: %s", description)
	}
}

func TestNewWeatherMetricsNotFoundError(t *testing.T) {
	testDate := "2024-01-15"
	err := domainError.NewWeatherMetricsNotFoundError(testDate)

	if err == nil {
		t.Error("NewWeatherMetricsNotFoundError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewWeatherMetricsNotFoundError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "no weather metrics found") {
		t.Errorf("NewWeatherMetricsNotFoundError() description should contain 'no weather metrics found', got: %s", description)
	}

	if !strings.Contains(description, testDate) {
		t.Errorf("NewWeatherMetricsNotFoundError() description should contain date '%s', got: %s", testDate, description)
	}
}

func TestNewInvalidCoordinatesError(t *testing.T) {
	lat := -91.0
	lng := 181.0
	err := domainError.NewInvalidCoordinatesError(lat, lng)

	if err == nil {
		t.Error("NewInvalidCoordinatesError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewInvalidCoordinatesError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "invalid coordinates") {
		t.Errorf("NewInvalidCoordinatesError() description should contain 'invalid coordinates', got: %s", description)
	}
}

func TestNewProcessingNotFoundError(t *testing.T) {
	processingID := "test-processing-id-123"
	err := domainError.NewProcessingNotFoundError(processingID)

	if err == nil {
		t.Error("NewProcessingNotFoundError() should return an error, got nil")
		return
	}

	baseErr, ok := err.(errs.BaseError)
	if !ok {
		t.Error("NewProcessingNotFoundError() should return a BaseError")
		return
	}

	description := baseErr.Description()
	if !strings.Contains(description, "processing") && !strings.Contains(description, "not found") {
		t.Errorf("NewProcessingNotFoundError() description should contain 'processing' and 'not found', got: %s", description)
	}
}
