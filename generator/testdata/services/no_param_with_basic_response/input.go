package no_param_with_basic_response

import (
	"context"
)

type Example interface {
	Method(ctx context.Context) (*string, error)
}
