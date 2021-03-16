package no_param_no_response

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
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

func (c *ExampleClient) Method(ctx context.Context) error {
	content := &client.DataContent{ContentType: "application/json"}
	_, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Method", "post", content)
	return err
}

func _Example_Method_Handler(srv Example) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		methodErr := srv.Method(ctx)
		if methodErr != nil {
			err = methodErr
			return
		}
		return
	}
}

type InvocationHandlerFunc func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)

func Register(s common.Service, srv Example) {
	s.AddServiceInvocationHandler("Method", _Example_Method_Handler(srv))
}

func NewExampleServer(address string, srv Example) (common.Service, error) {
	s, err := grpc.NewService(address)
	if err != nil {
		return nil, err
	}
	Register(s, srv)

	return s, nil
}
