package adapter

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context, event any) (any, error)
}
