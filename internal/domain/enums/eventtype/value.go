package eventtype

import (
	"strings"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/enums"
)

const (
	ForestCreated         enums.EventType = "forest_created"
	StartAnalysis         enums.EventType = "start_analysis"
	ReadyAnalysis         enums.EventType = "ready_analysis"
	Notify                enums.EventType = "notify"
	ForestCompletionCheck enums.EventType = "FOREST_COMPLETION_CHECK"
	Unknown               enums.EventType = "unknown"
)

func GetEventType(s string) enums.EventType {
	switch strings.ToLower(s) {
	case "forest_created":
		return ForestCreated
	case "start_analysis":
		return StartAnalysis
	case "ready_analysis":
		return ReadyAnalysis
	case "notify":
		return Notify
	case "forest_completion_check":
		return ForestCompletionCheck
	default:
		return Unknown
	}
}
