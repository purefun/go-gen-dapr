package echo

import "context"

type Message struct {
	Text string
}

type Service interface {
	Echo(ctx context.Context, in Message) (*Message, error)
}
