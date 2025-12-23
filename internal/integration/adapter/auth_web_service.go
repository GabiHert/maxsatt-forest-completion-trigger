package adapter

import (
	"context"
)

type AuthWebService interface {
	GetToken(ctx context.Context) (*string, error)
}
