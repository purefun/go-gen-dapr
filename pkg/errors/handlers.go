package errors

import (
	"context"
	"fmt"

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
			fmt.Println("Error:", err)
			if tracker, ok := err.(stackTracker); ok {
				for _, f := range tracker.StackTrace() {
					fmt.Printf("%+s:%d\n", f, f)
				}
			}
		}

		return
	}
}
