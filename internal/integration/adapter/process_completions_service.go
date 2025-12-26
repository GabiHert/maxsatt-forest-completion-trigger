package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/adapter"
)

// ProcessCompletionsService is the integration layer adapter for the process completions service.
// It embeds the application layer interface to ensure the integration layer
// can work with any implementation of the service.
type ProcessCompletionsService interface {
	adapter.ProcessCompletionsService
}
