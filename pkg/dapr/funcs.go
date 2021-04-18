package dapr

import (
	"context"

	"github.com/dapr/go-sdk/service/common"
)

type TopicHandlerFunc func(ctx context.Context, e *common.TopicEvent) (retry bool, err error)

type InvocationHandlerFunc func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
