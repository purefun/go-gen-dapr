package variadic

import (
	"context"
)

type Example interface {
	Method(ctx context.Context, ss ...string) error
}
