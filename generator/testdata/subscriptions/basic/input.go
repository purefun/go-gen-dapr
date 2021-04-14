package basic

import "context"

type UserRegisteredEvent struct {
	Email string
}

type Subscriber interface {
	SendActivationEmail(ctx context.Context, event UserRegisteredEvent) (bool, error)
}
