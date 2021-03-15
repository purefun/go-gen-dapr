// Code generated by go-gen-dapr. DO NOT EDIT.
// version: v0.1.0
// source:  github.com/purefun/go-gen-dapr/generator/testdata/example_service.Example

package example_service

import (
	"context"
	"encoding/json"
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

type InvocationHandlerFunc func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)

func (c *ExampleClient) Hello(ctx context.Context) {
	content := &client.DataContent{ContentType: "application/json"}
	c.cc.InvokeMethodWithContent(ctx, c.appID, "Hello", "post", content)
	return
}

func _Example_Hello_Handler(srv Example) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		srv.Hello(ctx)
		out = &common.Content{
			ContentType: "application/json",
		}
		return
	}
}

func Register(s common.Service, srv Example) {
	s.AddServiceInvocationHandler("Hello", _Example_Hello_Handler(srv))
}

func NewExampleServer(address string, srv Example) (common.Service, error) {
	s, err := grpc.NewService(address)
	if err != nil {
		return nil, err
	}
	Register(s, srv)

	return s, nil
}