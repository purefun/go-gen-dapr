package errors

import (
	"context"
	"log"

	"github.com/dapr/go-sdk/service/common"
	"github.com/pkg/errors"
	"github.com/purefun/go-gen-dapr/pkg/dapr"
)

type stackTracker interface {
	StackTrace() errors.StackTrace
}

func ServiceErrorHandler(handler dapr.InvocationHandlerFunc) dapr.InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out, err = handler(ctx, in)

		if err != nil {
			log.Println("service handler error:", err)
			if tracker, ok := err.(stackTracker); ok {
				for _, f := range tracker.StackTrace() {
					log.Printf("%+s:%d\n", f, f)
				}
			}
		}

		return
	}
}

func SubscriberErrorHandler(handler dapr.TopicHandlerFunc) dapr.TopicHandlerFunc {
	return func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		retry, err = handler(ctx, e)

		if err != nil {
			log.Printf("subscriber topic[%s] handler error: %s\n", e.Topic, err)
			if tracker, ok := err.(stackTracker); ok {
				for _, f := range tracker.StackTrace() {
					log.Printf("%+s:%d\n", f, f)
				}
			}
		}
		return
	}
}
