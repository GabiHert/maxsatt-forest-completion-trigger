package aws

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func SqsEventParser[T any](sqsMessage any, i *T) (map[string]any, int, error) {
	eventBytes, err := json.Marshal(sqsMessage)
	if err != nil {
		return nil, 0, errors.New("Error while trying to parse the event: " + err.Error())
	}

	var event *events.SQSEvent
	if err := json.Unmarshal(eventBytes, &event); err != nil {
		return nil, 0, errors.New("Error while trying to parse the event: " + err.Error())
	}

	var snsEntity events.SNSEntity
	var receiveCount int
	if err := json.Unmarshal([]byte(event.Records[0].Body), &snsEntity); err != nil {
		return nil, 0, errors.New("Error while trying to parse the sns entity: " + err.Error())
	}

	messageAttributes := map[string]any{}
	if snsEntity.Message != "" {
		if err := json.Unmarshal([]byte(snsEntity.Message), i); err != nil {
			return nil, 0, errors.New("Error while trying to parse the sns message: " + err.Error())
		}
		for key, value := range snsEntity.MessageAttributes {
			valueMap, ok := value.(map[string]any)
			if ok {
				messageAttributes[key] = valueMap["Value"]
			} else {
				messageAttributes[key] = value
			}
		}
	} else {
		if err := json.Unmarshal([]byte(event.Records[0].Body), i); err != nil {
			return nil, 0, errors.New("Error while trying to parse the sqs message: " + err.Error())
		}
		for key, value := range event.Records[0].MessageAttributes {
			messageAttributes[key] = *value.StringValue
		}
	}

	if value, ok := event.Records[0].Attributes["ApproximateReceiveCount"]; ok && value != "" {
		receiveCount, _ = strconv.Atoi(value)
	}

	return messageAttributes, receiveCount, nil
}

type DynamoDbImage string

const (
	NewImage DynamoDbImage = "NewImage"
	OldImage DynamoDbImage = "OldImage"
)

func DynamoDBStreamEventParser[T any](dynamoDbStreamMessage events.DynamoDBStreamRecord, imageType DynamoDbImage, object *T) error {
	if imageType == NewImage {
		if err := unmarshalStreamImage(dynamoDbStreamMessage.NewImage, object); err != nil {
			return errors.New("Error while trying to parse the new image: " + err.Error())
		}
	} else if imageType == OldImage {
		if err := unmarshalStreamImage(dynamoDbStreamMessage.OldImage, object); err != nil {
			return errors.New("Error while trying to parse the new image: " + err.Error())
		}
	}

	return nil
}

// UnmarshalStreamImage converts events.DynamoDBAttributeValue to struct
func unmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out any) error {
	dbAttrMap := make(map[string]types.AttributeValue)
	for k, v := range attribute {
		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return marshalErr
		}

		var dbAttr types.AttributeValue
		err := json.Unmarshal(bytes, &dbAttr)
		if err != nil {
			return err
		}
		dbAttrMap[k] = dbAttr
	}

	return attributevalue.UnmarshalMap(dbAttrMap, out)
}

// MarshalToStreamImage converts a struct to map[string]events.DynamoDBAttributeValue
func MarshalToStreamImage(object any) (map[string]events.DynamoDBAttributeValue, error) {
	dbAttrMap, err := attributevalue.MarshalMap(object)
	if err != nil {
		return nil, err
	}

	result := make(map[string]events.DynamoDBAttributeValue)
	for k, v := range dbAttrMap {
		bytes, marshalErr := json.Marshal(v)
		if marshalErr != nil {
			return nil, marshalErr
		}

		var filter map[string]any
		err = json.Unmarshal(bytes, &filter)
		if err != nil {
			return nil, err
		}

		filterNullValuesFromMap(filter)

		bytes, marshalErr = json.Marshal(filter)
		if marshalErr != nil {
			return nil, err
		}

		var eventAttr events.DynamoDBAttributeValue
		if err := json.Unmarshal(bytes, &eventAttr); err != nil {
			return nil, err
		}

		result[k] = eventAttr
	}

	return result, nil
}

func filterNullValuesFromMap(filter map[string]any) {
	for k, v := range filter {
		switch v.(type) {
		case map[string]any:
			filterNullValuesFromMap(v.(map[string]any))
		case []any:
			filterNullValuesFromList(v.([]any))
		case nil:
			delete(filter, k)
		}
	}
}

func filterNullValuesFromList(filter []any) {
	for i := 0; i < len(filter); i++ {
		switch filter[i].(type) {
		case map[string]any:
			filterNullValuesFromMap(filter[i].(map[string]any))
		case []any:
			filterNullValuesFromList(filter[i].([]any))
		case nil:
			filter = append(filter[:i], filter[i+1:]...)
		}
	}
}
