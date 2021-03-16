package external_types

import (
	"context"
	"github.com/purefun/go-gen-dapr/generator/testdata"
)

type Example interface {
	Method(ctx context.Context, a testdata.Input, b *testdata.Input) (*testdata.Output, error)
}
