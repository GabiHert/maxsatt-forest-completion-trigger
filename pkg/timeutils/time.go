package timeutils

import "time"

func ToDate(time time.Time) string {
	return time.UTC().Format("2006-01-02")
}

func FromDate(stringTime string) (*time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", stringTime)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

func ToIsoTimestamp(time time.Time) string {
	return time.UTC().Format("2006-01-02T15:04:05.000Z")
}

func FromIsoTimestamp(stringTime string) (*time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z", stringTime)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

func FromUtcTimestamp(stringTime string) (*time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", stringTime)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

func Now() time.Time {
	return GetTimeConfig().Now().UTC()
}

func CurrentDate() time.Time {
	date := Now()
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}
