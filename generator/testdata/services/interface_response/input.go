package interfaceresponse

import (
	"context"
)

type Example interface {
	Method(ctx context.Context) (interface{}, error)
}
