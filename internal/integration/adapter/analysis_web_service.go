package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
)

// AnalysisWebService defines the interface for retrieving analysis data from external API.
type AnalysisWebService interface {
	usecase.FetchFieldAnalysis
}
