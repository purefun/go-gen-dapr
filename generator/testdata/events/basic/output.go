package basic

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
)

type Pubsub struct {
	client     client.Client
	pubsubName string
}

func NewPubsub(c client.Client, name string) *Pubsub {
	return &Pubsub{client: c, pubsubName: name}
}

func (p *Pubsub) PublishUserRegisteredEvent(ctx context.Context, event UserRegisteredEvent) error {
	Data := struct {
		EventName string
		Event     UserRegisteredEvent
	}{
		EventName: "UserRegisteredEvent",
		Event:     event,
	}
	data, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.client.PublishEvent(ctx, p.pubsubName, "user", data)
}

func (p *Pubsub) PublishUserLoggedOutEvent(ctx context.Context, event UserLoggedOutEvent) error {
	Data := struct {
		EventName string
		Event     UserLoggedOutEvent
	}{
		EventName: "UserLoggedOutEvent",
		Event:     event,
	}
	data, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.client.PublishEvent(ctx, p.pubsubName, "user", data)
}
