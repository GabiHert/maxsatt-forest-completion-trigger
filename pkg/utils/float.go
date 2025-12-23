package utils

import (
	"strconv"
)

func ParseFloat(value string, defaultVal float32) float32 {
	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultVal
	}

	return float32(floatValue)
}
