package basic

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
)

type TopicHandlerFunc func(ctx context.Context, e *common.TopicEvent) (retry bool, err error)

func _SendActivationEmail_Handler(subscriber Subscriber) TopicHandlerFunc {
	return func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		var event UserRegisteredEvent
		err = json.Unmarshal(e.Data.([]byte), &event)
		if err != nil {
			return false, err
		}
		return subscriber.SendActivationEmail(ctx, event)
	}
}

func RegisterTopicHandlers(s common.Service, subscriber Subscriber, pubsubName string) {
	s.AddTopicEventHandler(
		&common.Subscription{
			PubsubName: pubsubName,
			Topic:      "UserRegisteredEvent",
		},
		_SendActivationEmail_Handler(subscriber),
	)
}
