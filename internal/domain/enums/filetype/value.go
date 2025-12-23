package filetype

import (
	"strings"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
)

const (
	Climate enums.FileType = "climate"
	Raw     enums.FileType = "raw"
	Clean   enums.FileType = "clean"
	Tif     enums.FileType = "tif"
	Unknown enums.FileType = "unknown"
)

func GetFileType(s string) enums.FileType {
	switch strings.ToLower(s) {
	case "climate":
		return Climate
	case "raw":
		return Raw
	case "clean":
		return Clean
	case "tif":
		return Tif
	default:
		return Unknown
	}
}
