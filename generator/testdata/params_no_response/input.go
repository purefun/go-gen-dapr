package params_no_response

import (
	"context"
)

type Example interface {
	ParamsNoResponse(ctx context.Context, a, b string) error
}
