package aws

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBHelperAdapter interface {
	BatchGetItem(ctx context.Context, key string, keyValues []string, table string, items any) error
	GetItem(ctx context.Context, key string, keyValue string, table string, item any) error
	GetItemWithRange(ctx context.Context, key string, keyValue string, rangeKey string, rangeValue string, table string, item any) error
	GetItems(ctx context.Context, table string, items any) error
	PutItem(ctx context.Context, table string, item any) error
	DeleteItem(ctx context.Context, key string, value string, table string) ([]byte, error)
	Query(ctx context.Context, input *dynamodb.QueryInput, items any) error
}

type DynamoDB interface {
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type dynamodbHelper struct {
	client   DynamoDB
	logger   loggerAdapter
	instance dynamodb.Client
}

func DynamoDBClient(region string) DynamoDB {
	return dynamodb.NewFromConfig(getConfig(region))
}

func DynamoDBHelper(client DynamoDB, logger loggerAdapter) DynamoDBHelperAdapter {
	return &dynamodbHelper{
		client: client,
		logger: logger,
	}
}

func (d *dynamodbHelper) BatchGetItem(ctx context.Context, key string, keyValues []string, table string, items any) error {
	d.logger.Debug(ctx, "Starting", key, keyValues, table)

	if table == "" {
		return errors.New("table name is required")
	}

	keys := make([]map[string]types.AttributeValue, 0, len(keyValues))
	for _, id := range keyValues {
		keys = append(keys, map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		})
	}

	response, err := d.client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			table: {
				Keys: keys,
			},
		},
	})
	if err != nil {
		return err
	}

	rawItems, ok := response.Responses[table]
	if !ok {
		return errors.New("error while trying to unmarshal item response")
	}

	err = attributevalue.UnmarshalListOfMaps(rawItems, &items)
	if err != nil {
		return errors.New("error while trying to unmarshal item response")
	}

	d.logger.Debug(ctx, "Finished", key, keyValues, table, items)
	return nil
}

func (d *dynamodbHelper) GetItem(ctx context.Context, key string, keyValue string, table string, item any) error {
	d.logger.Debug(ctx, "Starting", key, keyValue, table)

	if table == "" {
		return errors.New("table name is required")
	}

	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: keyValue},
		},
		TableName: &table,
	})
	if err != nil {
		return err
	}

	if response.Item == nil {
		return nil
	}

	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		return errors.New("error while trying to unmarshal item response")
	}

	d.logger.Debug(ctx, "Finished", key, keyValue, table, item)
	return nil
}

func (d *dynamodbHelper) GetItemWithRange(ctx context.Context,
	key string, keyValue string,
	rangeKey string, rangeValue string,
	table string,
	item any,
) error {
	d.logger.Debug(ctx, "Starting", key, keyValue, rangeKey, rangeValue, table)

	if table == "" {
		return errors.New("table name is required")
	}

	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			key:      &types.AttributeValueMemberS{Value: keyValue},
			rangeKey: &types.AttributeValueMemberS{Value: rangeValue},
		},
		TableName: &table,
	})
	if err != nil {
		return err
	}

	if response.Item == nil {
		return errors.New("item not found")
	}

	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		return errors.New("error while trying to unmarshal item response")
	}

	d.logger.Debug(ctx, "Finished", key, keyValue, table, item)
	return nil
}

func (d *dynamodbHelper) GetItems(ctx context.Context, table string, items any) error {
	d.logger.Debug(ctx, "Starting", table)

	if table == "" {
		return errors.New("table name is required")
	}

	response, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &table,
	})
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, items)
	if err != nil {
		return errors.New("error while trying to unmarshal item response")
	}

	d.logger.Debug(ctx, "Finished", table, items)
	return nil
}

func (d *dynamodbHelper) PutItem(ctx context.Context, table string, item any) error {
	d.logger.Debug(ctx, "Starting", table, item)

	if table == "" {
		return errors.New("table name is required")
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: &table,
	})
	if err != nil {
		return err
	}

	d.logger.Debug(ctx, "Finished", table, item)
	return nil
}

func (d *dynamodbHelper) DeleteItem(ctx context.Context, key string, value string, table string) ([]byte, error) {
	d.logger.Debug(ctx, "Starting", key, value, table)

	if table == "" {
		return nil, errors.New("table name is required")
	}

	result, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{Value: value},
		},
		TableName: &table,
	})
	if err != nil {
		return nil, err
	}

	var response map[string]any
	err = attributevalue.UnmarshalMap(result.Attributes, &response)
	if err != nil {
		return nil, errors.New("error while trying to unmarshal item response")
	}

	bytes, err := json.Marshal(response)

	d.logger.Debug(ctx, "Finished", key, value, table, bytes)
	return bytes, nil
}

func (d *dynamodbHelper) Query(ctx context.Context, input *dynamodb.QueryInput, items any) error {
	d.logger.Debug(ctx, "Starting", input, items)

	response, err := d.client.Query(ctx, input)
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, items)
	if err != nil {
		return errors.New("error while trying to unmarshal item response")
	}

	d.logger.Debug(ctx, "Finished", input, items)
	return nil
}
