package mock

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type SqsClient struct {
	Messages               map[string][]*string
	MessageAttributes      map[string][]*string
	MessageGroupId         map[string][]*string
	MessageDeduplicationId map[string][]*string
	id                     string
}

var sqsInstance *SqsClient
var sqsInit sync.Once

func NewSqsClient() *SqsClient {
	sqsInit.Do(
		func() {
			sqsInstance = &SqsClient{
				id:                     uuid.New().String(),
				Messages:               make(map[string][]*string),
				MessageAttributes:      make(map[string][]*string),
				MessageGroupId:         make(map[string][]*string),
				MessageDeduplicationId: make(map[string][]*string),
			}
		},
	)
	return sqsInstance
}

func (s *SqsClient) SendMessage(_ context.Context, input *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	if s.Messages[*input.QueueUrl] == nil {
		s.Messages[*input.QueueUrl] = make([]*string, 0)
	}
	s.Messages[*input.QueueUrl] = append(s.Messages[*input.QueueUrl], input.MessageBody)

	var messageAttributesString *string
	if input.MessageAttributes != nil {
		messageAttributes := make(map[string]*string)
		for key, value := range input.MessageAttributes {
			messageAttributes[key] = value.StringValue
		}
		messageAttributeBytes, err := json.Marshal(messageAttributes)
		if err != nil {
			return nil, err
		}

		value := string(messageAttributeBytes)
		messageAttributesString = &value
	}

	if s.MessageAttributes[*input.QueueUrl] == nil {
		s.MessageAttributes[*input.QueueUrl] = make([]*string, 0)
	}
	s.MessageAttributes[*input.QueueUrl] = append(s.MessageAttributes[*input.QueueUrl], messageAttributesString)

	if s.MessageGroupId[*input.QueueUrl] == nil {
		s.MessageGroupId[*input.QueueUrl] = make([]*string, 0)
	}
	s.MessageGroupId[*input.QueueUrl] = append(s.MessageGroupId[*input.QueueUrl], input.MessageGroupId)

	if s.MessageDeduplicationId[*input.QueueUrl] == nil {
		s.MessageDeduplicationId[*input.QueueUrl] = make([]*string, 0)
	}
	s.MessageDeduplicationId[*input.QueueUrl] = append(s.MessageDeduplicationId[*input.QueueUrl], input.MessageDeduplicationId)

	return nil, nil
}

func (s *SqsClient) Reset() {
	s.Messages = make(map[string][]*string)
	s.MessageAttributes = make(map[string][]*string)
	s.MessageGroupId = make(map[string][]*string)
	s.MessageDeduplicationId = make(map[string][]*string)
}
