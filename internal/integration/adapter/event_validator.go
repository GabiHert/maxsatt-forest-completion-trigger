package adapter

import "context"

type Validator[T any] interface {
	Validate(ctx context.Context, request *T) error
}
