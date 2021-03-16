package external_types

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/purefun/go-gen-dapr/generator/testdata"
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

func (c *ExampleClient) Method(ctx context.Context, a testdata.Input, b *testdata.Input) (*testdata.Output, error) {
	content := &client.DataContent{ContentType: "application/json"}
	params, encErr := json.Marshal(map[string]interface{}{
		"a": a,
		"b": b,
	})
	if encErr != nil {
		return encErr
	}
	content.Data = params
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Method", "post", content)
	var out *testdata.Output
	err := json.Unmarshal(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Example_Method_Handler(srv Example) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		type Params struct {
			A testdata.Input
			B *testdata.Input
		}
		var params Params
		decErr := json.Unmarshal(in.Data, &params)
		if decErr != nil {
			err = decErr
			return
		}
		resp, methodErr := srv.Method(ctx, params.A, params.B)
		if methodErr != nil {
			err = methodErr
			return
		}
		data, encErr := json.Marshal(resp)
		if encErr != nil {
			err = encErr
			return
		}
		out.Data = data
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
