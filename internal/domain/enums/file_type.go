package enums

import "strings"

type FileType string

func (f FileType) Equals(s string) bool {
	return strings.EqualFold(string(f), s)
}

func (f FileType) String() string {
	return strings.ToUpper(string(f))
}
