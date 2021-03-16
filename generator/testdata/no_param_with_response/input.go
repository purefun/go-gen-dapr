package no_param_with_response

import (
	"context"
)

type Example interface {
	NoParamWithResponse(ctx context.Context) (*string, error)
}
