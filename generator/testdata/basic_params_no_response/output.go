package basic_params_no_response

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

func (c *ExampleClient) Method(ctx context.Context, a string, b string) error {
	content := &client.DataContent{ContentType: "application/json"}
	params, encErr := json.Marshal([]interface{}{
		{Name: "a", Value: a},
		{Name: "b", Value: b},
	})
	if encErr != nil {
		return nil
	}
	content.Data = params
	_, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Method", "post", content)
	return err
}

func _Example_Method_Handler(srv Example) InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		var params *[]interface{}
		decErr := json.Unmarshal(in.Data, params)
		if decErr != nil {
			err = decErr
			return
		}
		_a, ok := (*params)[0].(string)
		if !ok {
			err = fmt.Errorf(`param "a" type is not "string"`)
			return
		}
		_b, ok := (*params)[1].(string)
		if !ok {
			err = fmt.Errorf(`param "b" type is not "string"`)
			return
		}
		methodErr := srv.Method(ctx, _a, _b)
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
