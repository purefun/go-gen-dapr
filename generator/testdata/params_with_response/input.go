package params_with_response

import (
	"context"
)

type Example interface {
	ParamsWithResponse(ctx context.Context, a, b string) (*string, error)
}
