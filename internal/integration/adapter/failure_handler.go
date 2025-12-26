package adapter

import (
	"context"
)

// FailureHandler defines the interface for handling processing failures.
type FailureHandler interface {
	Handle(ctx context.Context, err error) error
}
