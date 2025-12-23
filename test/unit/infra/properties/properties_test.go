package properties_test

import (
	"os"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
)

func TestProperties(t *testing.T) {
	_ = os.Setenv("SERVICE_NAME", "test-service")
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("FOREST_EVENTS_TOPIC_ARN", "arn:aws:sns:us-east-1:123456789012:test-topic")
	_ = os.Setenv("FOREST_SCHEDULE_QUEUE_URL", "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue")
	_ = os.Setenv("FOREST_SCHEDULE_INTERVAL_DAYS", "7")
	defer func() {
		_ = os.Unsetenv("SERVICE_NAME")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("FOREST_EVENTS_TOPIC_ARN")
		_ = os.Unsetenv("FOREST_SCHEDULE_QUEUE_URL")
		_ = os.Unsetenv("FOREST_SCHEDULE_INTERVAL_DAYS")
	}()

	props := properties.Properties()

	if props == nil {
		t.Fatal("Properties() returned nil")
	}

	if props.Application.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service", props.Application.ServiceName)
	}

	if props.Services.ForestEventsTopicArn != "arn:aws:sns:us-east-1:123456789012:test-topic" {
		t.Errorf("ForestEventsTopicArn = %v, want arn:aws:sns:us-east-1:123456789012:test-topic", props.Services.ForestEventsTopicArn)
	}

	if props.Services.ForestScheduleQueueUrl != "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue" {
		t.Errorf("ForestScheduleQueueUrl = %v, want https://sqs.us-east-1.amazonaws.com/123456789012/test-queue", props.Services.ForestScheduleQueueUrl)
	}

	if props.Services.ForestScheduleIntervalDays != 7 {
		t.Errorf("ForestScheduleIntervalDays = %v, want 7", props.Services.ForestScheduleIntervalDays)
	}
}

func TestProperties_DefaultValues(t *testing.T) {
	_ = os.Unsetenv("FOREST_SCHEDULE_INTERVAL_DAYS")
	_ = os.Unsetenv("FOREST_SCHEDULE_TABLE")

	props := properties.Properties()

	if props.Services.ForestScheduleIntervalDays != 5 {
		t.Errorf("ForestScheduleIntervalDays default = %v, want 5", props.Services.ForestScheduleIntervalDays)
	}

	if props.Services.ForestScheduleTable != "forests-analysis-schedules" {
		t.Errorf("ForestScheduleTable default = %v, want forests-analysis-schedules", props.Services.ForestScheduleTable)
	}
}

func TestProperties_InvalidIntervalDays(t *testing.T) {
	_ = os.Setenv("FOREST_SCHEDULE_INTERVAL_DAYS", "invalid")
	defer func() {
		_ = os.Unsetenv("FOREST_SCHEDULE_INTERVAL_DAYS")
	}()

	props := properties.Properties()

	if props.Services.ForestScheduleIntervalDays != 5 {
		t.Errorf("ForestScheduleIntervalDays with invalid value = %v, want 5 (default)", props.Services.ForestScheduleIntervalDays)
	}
}

func TestProperties_EmptyTableNameUsesDefault(t *testing.T) {
	_ = os.Setenv("FOREST_SCHEDULE_TABLE", "")
	defer func() {
		_ = os.Unsetenv("FOREST_SCHEDULE_TABLE")
	}()

	props := properties.Properties()

	if props.Services.ForestScheduleTable != "forests-analysis-schedules" {
		t.Errorf("ForestScheduleTable with empty value = %v, want forests-analysis-schedules (default)", props.Services.ForestScheduleTable)
	}
}

func TestProperties_CustomTableName(t *testing.T) {
	_ = os.Setenv("FOREST_SCHEDULE_TABLE", "custom-table-name")
	defer func() {
		_ = os.Unsetenv("FOREST_SCHEDULE_TABLE")
	}()

	props := properties.Properties()

	if props.Services.ForestScheduleTable != "custom-table-name" {
		t.Errorf("ForestScheduleTable = %v, want custom-table-name", props.Services.ForestScheduleTable)
	}
}
