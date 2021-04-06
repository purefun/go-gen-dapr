package external_types

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"github.com/purefun/go-gen-dapr/generator/testdata/subscriptions"
)

type TopicHandlerFunc func(ctx context.Context, e *common.TopicEvent) (retry bool, err error)

func _SendActivationEmail_Handler(subs Subscriptions) TopicHandlerFunc {
	return func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		var event subscriptions.UserRegisteredEvent
		err = json.Unmarshal(e.Data.([]byte), &event)
		if err != nil {
			return false, err
		}
		return subs.SendActivationEmail(ctx, event)
	}
}

func RegisterTopicHandlers(s common.Service, subs Subscriptions, pubsubName string) {
	s.AddTopicEventHandler(
		&common.Subscription{
			PubsubName: pubsubName,
			Topic:      "UserRegisteredEvent",
		},
		_SendActivationEmail_Handler(subs),
	)
}
