package main

import "context"

type Message struct {
	Text string
}

type Service interface {
	Echo(ctx context.Context) (*string, error)
	Hello(ctx context.Context, in Message) (*Message, error)
}
