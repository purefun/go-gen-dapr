package external_types

import (
	"context"
	"github.com/purefun/go-gen-dapr/generator/testdata/services"
)

type Example interface {
	Method(ctx context.Context, a services.Input, b *services.Input) (*services.Output, error)
}
