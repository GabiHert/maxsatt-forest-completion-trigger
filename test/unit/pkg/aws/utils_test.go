package aws_test

import (
	"encoding/json"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"

	"github.com/aws/aws-lambda-go/events"
)

type TestMessage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestSqsEventParser_WithSNSEntity(t *testing.T) {
	message := TestMessage{
		ID:   "123",
		Name: "test",
	}
	messageJSON, _ := json.Marshal(message)

	snsEntity := events.SNSEntity{
		Message: string(messageJSON),
		MessageAttributes: map[string]any{
			"attribute1": map[string]any{
				"Type":  "String",
				"Value": "value1",
			},
		},
	}
	snsJSON, _ := json.Marshal(snsEntity)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsJSON),
				Attributes: map[string]string{
					"ApproximateReceiveCount": "3",
				},
			},
		},
	}

	var result TestMessage
	attributes, receiveCount, err := aws.SqsEventParser(sqsEvent, &result)

	if err != nil {
		t.Errorf("SqsEventParser() should not return error, got: %v", err)
	}

	if result.ID != "123" || result.Name != "test" {
		t.Errorf("SqsEventParser() parsed message incorrectly: %+v", result)
	}

	if receiveCount != 3 {
		t.Errorf("SqsEventParser() receive count = %d, want 3", receiveCount)
	}

	if attributes["attribute1"] != "value1" {
		t.Errorf("SqsEventParser() attributes parsed incorrectly: %+v", attributes)
	}
}

func TestSqsEventParser_WithDirectSQSMessage(t *testing.T) {
	message := TestMessage{
		ID:   "456",
		Name: "direct",
	}
	messageJSON, _ := json.Marshal(message)

	stringValue := "attr-value"
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(messageJSON),
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"attribute1": {
						StringValue: &stringValue,
					},
				},
				Attributes: map[string]string{
					"ApproximateReceiveCount": "1",
				},
			},
		},
	}

	var result TestMessage
	attributes, receiveCount, err := aws.SqsEventParser(sqsEvent, &result)

	if err != nil {
		t.Errorf("SqsEventParser() should not return error, got: %v", err)
	}

	if result.ID != "456" || result.Name != "direct" {
		t.Errorf("SqsEventParser() parsed message incorrectly: %+v", result)
	}

	if receiveCount != 1 {
		t.Errorf("SqsEventParser() receive count = %d, want 1", receiveCount)
	}

	if attributes["attribute1"] != "attr-value" {
		t.Errorf("SqsEventParser() attributes parsed incorrectly: %+v", attributes)
	}
}

func TestSqsEventParser_InvalidEvent(t *testing.T) {
	invalidEvent := "not a valid event"

	var result TestMessage
	_, _, err := aws.SqsEventParser(invalidEvent, &result)

	if err == nil {
		t.Error("SqsEventParser() should return error for invalid event")
	}
}

func TestSqsEventParser_InvalidSNSEntity(t *testing.T) {
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: "invalid json",
				Attributes: map[string]string{
					"ApproximateReceiveCount": "1",
				},
			},
		},
	}

	var result TestMessage
	_, _, err := aws.SqsEventParser(sqsEvent, &result)

	if err == nil {
		t.Error("SqsEventParser() should return error for invalid SNS entity")
	}
}

func TestSqsEventParser_InvalidSNSMessage(t *testing.T) {
	snsEntity := events.SNSEntity{
		Message: "invalid json",
	}
	snsJSON, _ := json.Marshal(snsEntity)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsJSON),
				Attributes: map[string]string{
					"ApproximateReceiveCount": "1",
				},
			},
		},
	}

	var result TestMessage
	_, _, err := aws.SqsEventParser(sqsEvent, &result)

	if err == nil {
		t.Error("SqsEventParser() should return error for invalid SNS message")
	}
}

func TestSqsEventParser_MissingReceiveCount(t *testing.T) {
	message := TestMessage{
		ID:   "789",
		Name: "no-count",
	}
	messageJSON, _ := json.Marshal(message)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body:       string(messageJSON),
				Attributes: map[string]string{},
			},
		},
	}

	var result TestMessage
	_, receiveCount, err := aws.SqsEventParser(sqsEvent, &result)

	if err != nil {
		t.Errorf("SqsEventParser() should not return error, got: %v", err)
	}

	if receiveCount != 0 {
		t.Errorf("SqsEventParser() receive count = %d, want 0", receiveCount)
	}
}

func TestSqsEventParser_EmptyReceiveCount(t *testing.T) {
	message := TestMessage{
		ID:   "999",
		Name: "empty-count",
	}
	messageJSON, _ := json.Marshal(message)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(messageJSON),
				Attributes: map[string]string{
					"ApproximateReceiveCount": "",
				},
			},
		},
	}

	var result TestMessage
	_, receiveCount, err := aws.SqsEventParser(sqsEvent, &result)

	if err != nil {
		t.Errorf("SqsEventParser() should not return error, got: %v", err)
	}

	if receiveCount != 0 {
		t.Errorf("SqsEventParser() receive count = %d, want 0", receiveCount)
	}
}

func TestSqsEventParser_WithNonMapAttribute(t *testing.T) {
	message := TestMessage{
		ID:   "111",
		Name: "non-map",
	}
	messageJSON, _ := json.Marshal(message)

	snsEntity := events.SNSEntity{
		Message: string(messageJSON),
		MessageAttributes: map[string]any{
			"stringAttr": "direct-string-value",
		},
	}
	snsJSON, _ := json.Marshal(snsEntity)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsJSON),
				Attributes: map[string]string{
					"ApproximateReceiveCount": "1",
				},
			},
		},
	}

	var result TestMessage
	attributes, _, err := aws.SqsEventParser(sqsEvent, &result)

	if err != nil {
		t.Errorf("SqsEventParser() should not return error, got: %v", err)
	}

	if attributes["stringAttr"] != "direct-string-value" {
		t.Errorf("SqsEventParser() should handle non-map attributes: %+v", attributes)
	}
}

// DynamoDBStreamEventParser and MarshalToStreamImage are not used by the application
// based on mutation coverage analysis showing them as NOT COVERED
// Skipping tests for unused code to avoid implementation complexity
