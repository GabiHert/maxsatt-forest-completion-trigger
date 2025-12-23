package adapter

import "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"

type DeltaDatasetPersistenceAdapter interface {
	usecase.ProcessParquetInChunks
}
