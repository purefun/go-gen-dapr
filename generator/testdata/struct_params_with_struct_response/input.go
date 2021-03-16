package struct_params_with_struct_response

import (
	"context"
)

type Input struct {
	Name string
}

type Output struct {
	ID   string
	Name string
}

type Example interface {
	Method(ctx context.Context, a Input, b *Input) (*Output, error)
}
