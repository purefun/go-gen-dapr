package params_with_response

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

func (c *ExampleClient) ParamsWithResponse(ctx context.Context, a string, b string) (*string, error) {
	content := &client.DataContent{ContentType: "application/json"}
	params, encErr := json.Marshal([]interface{}{
		{Name: "a", Value: a},
		{Name: "b", Value: b},
	})
	if encErr != nil {
		return nil
	}
	content.Data = params
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "ParamsWithResponse", "post", content)
	var out *string
	err := json.Unmarshal(resp, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Example_ParamsWithResponse_Handler(srv Example) InvocationHandlerFunc {
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
		resp, methodErr := srv.ParamsWithResponse(ctx, _a, _b)
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
	s.AddServiceInvocationHandler("ParamsWithResponse", _Example_ParamsWithResponse_Handler(srv))
}

func NewExampleServer(address string, srv Example) (common.Service, error) {
	s, err := grpc.NewService(address)
	if err != nil {
		return nil, err
	}
	Register(s, srv)

	return s, nil
}
