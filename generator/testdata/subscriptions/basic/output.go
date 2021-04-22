package basic

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"github.com/purefun/go-gen-dapr/pkg/dapr"
	errorHandlers "github.com/purefun/go-gen-dapr/pkg/errors"
)

func _SendActivationEmail_Handler(subscriber Subscriber) dapr.TopicHandlerFunc {
	return func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		var event UserRegisteredEvent
		err = json.Unmarshal(e.Data.([]byte), &event)
		if err != nil {
			return false, err
		}
		return subscriber.SendActivationEmail(ctx, event)
	}
}

func RegisterSubscriber(s common.Service, subscriber Subscriber, pubsubName string) {
	s.AddTopicEventHandler(
		&common.Subscription{
			PubsubName: pubsubName,
			Topic:      "UserRegisteredEvent",
		},
		errorHandlers.SubscriberErrorHandler(_SendActivationEmail_Handler(subscriber)),
	)
}
