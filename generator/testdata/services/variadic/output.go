package variadic

import (
	"context"
	"encoding/json"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/pkg/errors"
	"github.com/purefun/go-gen-dapr/pkg/dapr"
	errorHandlers "github.com/purefun/go-gen-dapr/pkg/errors"
)

type ExampleClient struct {
	cc    client.Client
	appID string
}

func NewExampleClient(appID string) (*ExampleClient, error) {
	cc, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	return &ExampleClient{cc, appID}, nil
}

func (c *ExampleClient) Method(ctx context.Context, ss ...string) error {
	_content := &client.DataContent{ContentType: "application/json"}
	params, encErr := json.Marshal(map[string]interface{}{
		"ss": ss,
	})
	if encErr != nil {
		return errors.WithStack(encErr)
	}
	_content.Data = params
	_, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Method", "post", _content)
	return errors.WithStack(err)
}

func _Example_Method_Handler(srv Example) dapr.InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		type Params struct {
			Ss []string
		}
		var params Params
		decErr := json.Unmarshal(in.Data, &params)
		if decErr != nil {
			err = errors.WithStack(decErr)
			return
		}
		methodErr := srv.Method(ctx, params.Ss...)
		if methodErr != nil {
			err = errors.WithStack(methodErr)
			return
		}
		return
	}
}

func RegisterService(s common.Service, srv Example) {
	s.AddServiceInvocationHandler("Method", errorHandlers.ServiceErrorHandler(_Example_Method_Handler(srv)))
}

func NewExampleServer(address string) (common.Service, error) {
	s, err := grpc.NewService(address)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s, nil
}

func MustNewExampleServer(address string) common.Service {
	svc, err := NewExampleServer(address)
	if err != nil {
		panic(err)
	}
	return svc
}
