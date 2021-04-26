package slice_response

import (
	"context"
)

type Example interface {
	Method(ctx context.Context) ([]*string, error)
}
