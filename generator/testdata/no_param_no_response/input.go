package no_param_no_response

import (
	"context"
)

type Example interface {
	Method(ctx context.Context) error
}
