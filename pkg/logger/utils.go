package logger

import (
	"reflect"
	"strconv"
)

func parseBool(string string) bool {
	boolValue, err := strconv.ParseBool(string)
	if err != nil {
		return false
	}
	return boolValue
}

func entityName(entity any) string {
	return reflect.TypeOf(entity).Name()
}
