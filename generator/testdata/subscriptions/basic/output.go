package basic

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
)

type TopicHandlerFunc func(ctx context.Context, e *common.TopicEvent) (retry bool, err error)

func _UserRegisteredEvent_Handler(subs Subscriptions) TopicHandlerFunc {
	return func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		var event UserRegisteredEvent
		err = json.Unmarshal(e.Data, &event)
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
		_UserRegisteredEvent_Handler(subs),
	)
}
