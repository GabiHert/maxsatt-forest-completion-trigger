package enums

import "strings"

type EventType string

func (e EventType) Equals(s string) bool {
	return strings.EqualFold(string(e), s)
}

func (e EventType) String() string {
	return strings.ToUpper(string(e))
}
