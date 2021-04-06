package no_ctx_param

import "context"

type InvalidService interface {
	Hello() string
	Echo(text string) string
}
