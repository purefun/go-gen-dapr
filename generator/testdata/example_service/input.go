package example_service

import (
	"context"
)

type Example interface {
	NoResponse(ctx context.Context) error
	WithResponse(ctx context.Context) (*string, error)
}
