package no_param_no_response

import (
	"context"
)

type Example interface {
	NoParamNoResponse(ctx context.Context) error
}
