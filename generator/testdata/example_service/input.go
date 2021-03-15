package example_service

import (
	"context"
)

type Example interface {
	Hello(ctx context.Context)
}
