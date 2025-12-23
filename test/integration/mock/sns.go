package mock

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
)

type SnsClient struct {
	NumberOfMessagesSent map[string]int
	Messages             map[string]*string
	MessageAttributes    map[string]*string
	id                   string
}

var snsInstance *SnsClient
var snsInit sync.Once

func NewSnsClient() *SnsClient {
	snsInit.Do(
		func() {
			snsInstance = &SnsClient{
				id:                   uuid.New().String(),
				NumberOfMessagesSent: make(map[string]int),
				Messages:             make(map[string]*string),
				MessageAttributes:    make(map[string]*string),
			}
		},
	)
	return snsInstance
}

func (s *SnsClient) Publish(_ context.Context, input *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	s.Messages[*input.TopicArn] = input.Message
	s.NumberOfMessagesSent[*input.TopicArn]++

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

	s.MessageAttributes[*input.TopicArn] = messageAttributesString

	return nil, nil
}

func (s *SnsClient) Reset() {
	s.NumberOfMessagesSent = make(map[string]int)
	s.Messages = make(map[string]*string)
	s.MessageAttributes = make(map[string]*string)
}
