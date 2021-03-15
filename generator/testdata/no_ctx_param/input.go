package no_ctx_param

import "context"

type Service interface {
	Hello() string
	Echo(text string) string
}
