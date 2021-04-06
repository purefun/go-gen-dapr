package basic_params_with_basic_response

import (
	"context"
)

type Example interface {
	Method(ctx context.Context, a, b string) (*string, error)
}
