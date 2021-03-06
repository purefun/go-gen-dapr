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
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.PublishEvent(ctx, p.pubsubName, "UserRegisteredEvent", data)
}

func (p *Pubsub) PublishUserLoggedOutEvent(ctx context.Context, event UserLoggedOutEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.PublishEvent(ctx, p.pubsubName, "UserLoggedOutEvent", data)
}
