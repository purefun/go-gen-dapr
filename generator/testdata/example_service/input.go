package example_service

import (
	"context"
)

type Example interface {
	NoResponse(ctx context.Context) error
	WithResponse(ctx context.Context) (*string, error)
	Param(ctx context.Context, in string) error
	Params(ctx context.Context, a, b, c string) error
}
