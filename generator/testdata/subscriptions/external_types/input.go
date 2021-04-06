package external_types

import (
	"github.com/purefun/go-gen-dapr/generator/testdata/subscriptions"
)

type Subscriptions interface {
	SendActivationEmail(ctx context.Context, event subscriptions.UserRegisteredEvent) (bool, error)
}
