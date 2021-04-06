package service

import "context"

type Message struct {
	Text string
}

type Service interface {
	Echo(ctx context.Context) (*string, error)
	Hello(ctx context.Context, in Message) (*Message, error)
}

type UserRegisteredEvent struct {
	Email string
}

type Subscriptions interface {
	SendEmail(ctx context.Context, event UserRegisteredEvent) (bool, error)
}
