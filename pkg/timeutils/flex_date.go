package timeutils

import (
	"encoding/json"
	"time"
)

// FlexDate is a custom type that can unmarshal both RFC3339 and YYYY-MM-DD date formats.
type FlexDate struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface to handle multiple date formats.
func (fd *FlexDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		fd.Time = t
		return nil
	}

	if t, err := time.Parse("2006-01-02", s); err == nil {
		fd.Time = t
		return nil
	}

	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		fd.Time = t
		return nil
	}

	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (fd *FlexDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.Time)
}
