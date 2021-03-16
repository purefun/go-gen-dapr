package basic_params_no_response

import (
	"context"
)

type Example interface {
	Method(ctx context.Context, a, b string) error
}
