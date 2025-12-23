package mock

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoDbClient struct {
	items     map[string][]map[string]any
	tableKeys map[string]string
	id        string
}

var (
	dynamoDbMock *DynamoDbClient
	dynamoDbOnce sync.Once
)

func NewDynamoDbClient() *DynamoDbClient {
	dynamoDbOnce.Do(func() {
		dynamoDbMock = &DynamoDbClient{
			id:        uuid.NewString(),
			items:     map[string][]map[string]any{},
			tableKeys: map[string]string{},
		}
	})
	return dynamoDbMock
}

func (d *DynamoDbClient) AddTable(tableName, tableKey string) {
	d.tableKeys[tableName] = tableKey
}

func (d *DynamoDbClient) AddItem(item map[string]any, tableName string) {
	key := d.tableKeys[tableName]
	keyValue := item[key].(string)

	for i, oldItem := range d.items[tableName] {
		if oldItem[key] == keyValue {
			d.items[tableName] = append(d.items[tableName][:i], d.items[tableName][i+1:]...)
			break
		}
	}
	d.items[tableName] = append(d.items[tableName], item)
}

func (d *DynamoDbClient) GetAllItems(tableName string) []map[string]any {
	return d.items[tableName]
}

func (d *DynamoDbClient) Reset() *DynamoDbClient {
	d.items = map[string][]map[string]any{}
	d.tableKeys = map[string]string{}
	return d
}

func (d *DynamoDbClient) ResetTable(tableName string) error {
	if _, exists := d.tableKeys[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	d.items[tableName] = []map[string]any{}

	if len(d.items[tableName]) != 0 {
		return fmt.Errorf("failed to clear table %s: still contains %d items", tableName, len(d.items[tableName]))
	}

	return nil
}

// -----------------------------
// CRUD Methods
// -----------------------------

func (d *DynamoDbClient) GetItem(_ context.Context, input *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	key := map[string]string{}
	_ = attributevalue.UnmarshalMap(input.Key, &key)

	item := d.findMatchingItem(*input.TableName, key)
	if item == nil {
		return &dynamodb.GetItemOutput{}, nil
	}

	av, _ := attributevalue.MarshalMap(item)
	return &dynamodb.GetItemOutput{Item: av}, nil
}

func (d *DynamoDbClient) PutItem(_ context.Context, input *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	item := map[string]any{}
	_ = attributevalue.UnmarshalMap(input.Item, &item)

	tableName := *input.TableName
	keyField := d.tableKeys[tableName]

	keyValue, ok := item[keyField].(string)
	if !ok {
		return nil, fmt.Errorf("missing primary key: %s", keyField)
	}

	for i, oldItem := range d.items[tableName] {
		if oldItem[keyField] == keyValue {
			d.items[tableName] = append(d.items[tableName][:i], d.items[tableName][i+1:]...)
			break
		}
	}
	d.items[tableName] = append(d.items[tableName], item)
	return &dynamodb.PutItemOutput{}, nil
}

func (d *DynamoDbClient) DeleteItem(_ context.Context, input *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	key := map[string]string{}
	_ = attributevalue.UnmarshalMap(input.Key, &key)

	table := *input.TableName
	for i, item := range d.items[table] {
		if matches(item, key) {
			d.items[table] = append(d.items[table][:i], d.items[table][i+1:]...)
			av, _ := attributevalue.MarshalMap(item)
			return &dynamodb.DeleteItemOutput{Attributes: av}, nil
		}
	}
	return &dynamodb.DeleteItemOutput{}, nil
}

func (d *DynamoDbClient) Scan(_ context.Context, input *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	var items []map[string]types.AttributeValue
	for _, item := range d.items[*input.TableName] {
		av, _ := attributevalue.MarshalMap(item)
		items = append(items, av)
	}
	return &dynamodb.ScanOutput{Items: items}, nil
}

func (d *DynamoDbClient) Query(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	keyFilters := d.getKeyFilters(input)
	filterFilters := d.getFilterFilters(input)

	results := d.filterItems(*input.TableName, keyFilters, filterFilters)

	startIndex := d.getStartIndex(input.ExclusiveStartKey, results)

	limit := d.getLimit(input.Limit, len(results)-startIndex)

	pagedResults := results[startIndex : startIndex+limit]

	marshaled := d.marshalItems(pagedResults)

	output := &dynamodb.QueryOutput{
		Items: marshaled,
	}

	if startIndex+limit < len(results) {
		output.LastEvaluatedKey = d.getLastEvaluatedKey(*input.TableName, results[startIndex+limit-1])
	}

	return output, nil
}

func (d *DynamoDbClient) BatchGetItem(_ context.Context, input *dynamodb.BatchGetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	responses := make(map[string][]map[string]types.AttributeValue)

	for tableName, keysAndAttrs := range input.RequestItems {
		var tableResponses []map[string]types.AttributeValue
		for _, keyMap := range keysAndAttrs.Keys {
			key := map[string]string{}
			_ = attributevalue.UnmarshalMap(keyMap, &key)

			item := d.findMatchingItem(tableName, key)
			if item != nil {
				av, _ := attributevalue.MarshalMap(item)
				tableResponses = append(tableResponses, av)
			}
		}
		responses[tableName] = tableResponses
	}

	return &dynamodb.BatchGetItemOutput{
		Responses: responses,
	}, nil
}

// -----------------------------
// Helpers
// -----------------------------

func (d *DynamoDbClient) getKeyFilters(input *dynamodb.QueryInput) []filter {
	if input.KeyConditionExpression != nil {
		return parseKeyCondition(*input.KeyConditionExpression, input.ExpressionAttributeValues)
	}
	return nil
}

func (d *DynamoDbClient) getFilterFilters(input *dynamodb.QueryInput) []filter {
	if input.FilterExpression != nil {
		return parseFilterExpression(input.FilterExpression, input.ExpressionAttributeValues)
	}
	return nil
}

func (d *DynamoDbClient) filterItems(tableName string, keyFilters, filterFilters []filter) []map[string]any {
	var results []map[string]any
	for _, item := range d.items[tableName] {
		if matchesAll(item, keyFilters) && matchesAll(item, filterFilters) {
			results = append(results, item)
		}
	}
	return results
}

func (d *DynamoDbClient) getStartIndex(exclusiveStartKey map[string]types.AttributeValue, results []map[string]any) int {
	if exclusiveStartKey == nil {
		return 0
	}

	startKey := map[string]any{}
	_ = attributevalue.UnmarshalMap(exclusiveStartKey, &startKey)

	for i, item := range results {
		if matches(item, toStringMap(startKey)) {
			return i + 1
		}
	}
	return 0
}

func (d *DynamoDbClient) getLimit(limit *int32, remaining int) int {
	if limit != nil && *limit < int32(remaining) {
		return int(*limit)
	}
	return remaining
}

func (d *DynamoDbClient) marshalItems(items []map[string]any) []map[string]types.AttributeValue {
	var marshaled []map[string]types.AttributeValue
	for _, item := range items {
		av, _ := attributevalue.MarshalMap(item)
		marshaled = append(marshaled, av)
	}
	return marshaled
}

func (d *DynamoDbClient) getLastEvaluatedKey(tableName string, lastItem map[string]any) map[string]types.AttributeValue {
	lastKeyAV, _ := attributevalue.MarshalMap(map[string]any{
		d.tableKeys[tableName]: lastItem[d.tableKeys[tableName]],
	})
	return lastKeyAV
}

func (d *DynamoDbClient) findMatchingItem(table string, match map[string]string) map[string]any {
	for _, item := range d.items[table] {
		if matches(item, match) {
			return item
		}
	}
	return nil
}

func matches(item map[string]any, match map[string]string) bool {
	for k, v := range match {
		if fmt.Sprintf("%v", item[k]) != v {
			return false
		}
	}
	return true
}

type filter struct {
	Value    any
	Field    string
	Operator string
}

func parseKeyCondition(expr string, values map[string]types.AttributeValue) []filter {
	if expr == "" {
		return nil
	}
	var filters []filter
	parts := strings.Split(expr, "AND")
	for _, part := range parts {
		for _, op := range []string{"=", "<=", ">=", "<", ">"} {
			if strings.Contains(part, op) {
				tokens := strings.Split(part, op)
				field := strings.TrimSpace(tokens[0])
				valKey := strings.TrimSpace(tokens[1])
				val := extractValue(values[valKey])
				filters = append(filters, filter{Field: field, Operator: op, Value: val})
				break
			}
		}
	}
	return filters
}

func parseFilterExpression(expr *string, values map[string]types.AttributeValue) []filter {
	if expr == nil {
		return nil
	}
	var filters []filter
	parts := strings.Split(*expr, "AND")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "attribute_not_exists(") {
			field := strings.TrimPrefix(part, "attribute_not_exists(")
			field = strings.TrimSuffix(field, ")")
			field = strings.TrimSpace(field)
			filters = append(filters, filter{Field: field, Operator: "attribute_not_exists"})
			continue
		}

		if strings.HasPrefix(part, "contains(") {
			inner := strings.TrimPrefix(part, "contains(")
			inner = strings.TrimSuffix(inner, ")")
			args := strings.Split(inner, ",")
			if len(args) != 2 {
				continue
			}
			field := strings.TrimSpace(args[0])
			valKey := strings.TrimSpace(args[1])
			val := extractValue(values[valKey])
			filters = append(filters, filter{Field: field, Operator: "contains", Value: val})
			continue
		}

		for _, op := range []string{"=", "<=", ">=", "<", ">"} {
			if strings.Contains(part, op) {
				tokens := strings.Split(part, op)
				field := strings.TrimSpace(tokens[0])
				valKey := strings.TrimSpace(tokens[1])
				val := extractValue(values[valKey])
				filters = append(filters, filter{Field: field, Operator: op, Value: val})
				break
			}
		}
	}
	return filters
}

func matchesAll(item map[string]any, filters []filter) bool {
	for _, f := range filters {
		v := item[f.Field]

		switch f.Operator {
		case "=":
			if v != f.Value {
				return false
			}
		case "<":
			if !compareTime(v, f.Value, "<") {
				return false
			}
		case ">":
			if !compareTime(v, f.Value, ">") {
				return false
			}
		case "<=":
			if !compareTime(v, f.Value, "<=") {
				return false
			}
		case ">=":
			if !compareTime(v, f.Value, ">=") {
				return false
			}
		case "contains":
			strVal := fmt.Sprintf("%v", v)
			substr := fmt.Sprintf("%v", f.Value)
			if !strings.Contains(strVal, substr) {
				return false
			}
		case "attribute_not_exists":
			if _, exists := item[f.Field]; exists {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func toStringMap(m map[string]any) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

func compareTime(a, b any, op string) bool {
	at, aok := toTime(a)
	bt, bok := toTime(b)
	if aok && bok {
		switch op {
		case "<":
			return at.Before(bt)
		case ">":
			return at.After(bt)
		case "<=":
			return at.Before(bt) || at.Equal(bt)
		case ">=":
			return at.After(bt) || at.Equal(bt)
		}
	}

	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		switch op {
		case "<":
			return af < bf
		case ">":
			return af > bf
		case "<=":
			return af <= bf
		case ">=":
			return af >= bf
		}
	}

	return false
}

func extractValue(attr types.AttributeValue) any {
	switch v := attr.(type) {
	case *types.AttributeValueMemberS:
		return v.Value
	case *types.AttributeValueMemberN:
		if f, err := strconv.ParseFloat(v.Value, 64); err == nil {
			return f
		}
		return v.Value
	default:
		return nil
	}
}

func toTime(v any) (time.Time, bool) {
	switch t := v.(type) {
	case time.Time:
		return t, true
	case string:
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func toFloat64(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
