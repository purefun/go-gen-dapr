package basic

import "context"

type UserRegisteredEvent struct {
	Email string
}

type Subscriptions interface {
	SendActivationEmail(ctx context.Context, event UserRegisteredEvent) (bool, error)
}
