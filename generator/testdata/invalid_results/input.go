package no_error_result

import "context"

type Service interface {
	Hello(ctx context.Context)
	Hello1(ctx context.Context) string
	Hello2(ctx context.Context) *string
	Hello3(ctx context.Context) (string, error)
	Hello4(ctx context.Context) (*string, *string, error)
}
