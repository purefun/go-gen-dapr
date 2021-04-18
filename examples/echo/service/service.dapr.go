// Code generated by go-gen-dapr. DO NOT EDIT.
// version: v0.5.0
// source:  github.com/purefun/go-gen-dapr/examples/echo/service.Service

package service

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/purefun/go-gen-dapr/pkg/dapr"
	"github.com/purefun/go-gen-dapr/pkg/errors"
)

type ServiceClient struct {
	cc    client.Client
	appID string
}

func NewServiceClient(appID string) (*ServiceClient, error) {
	cc, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	return &ServiceClient{cc, appID}, nil
}

func (c *ServiceClient) Echo(ctx context.Context) (*string, error) {
	content := &client.DataContent{ContentType: "application/json"}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Echo", "post", content)
	if err != nil {
		return nil, err
	}
	if string(resp) == "null" {
		return nil, nil
	}
	var out string
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func _Service_Echo_Handler(srv Service) dapr.InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		resp, methodErr := srv.Echo(ctx)
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

func (c *ServiceClient) Hello(ctx context.Context, in Message) (*Message, error) {
	content := &client.DataContent{ContentType: "application/json"}
	params, encErr := json.Marshal(map[string]interface{}{
		"in": in,
	})
	if encErr != nil {
		return nil, encErr
	}
	content.Data = params
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "Hello", "post", content)
	if err != nil {
		return nil, err
	}
	if string(resp) == "null" {
		return nil, nil
	}
	var out Message
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func _Service_Hello_Handler(srv Service) dapr.InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		type Params struct {
			In Message
		}
		var params Params
		decErr := json.Unmarshal(in.Data, &params)
		if decErr != nil {
			err = decErr
			return
		}
		resp, methodErr := srv.Hello(ctx, params.In)
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

func (c *ServiceClient) SomethingWrong(ctx context.Context) (*Message, error) {
	content := &client.DataContent{ContentType: "application/json"}
	resp, err := c.cc.InvokeMethodWithContent(ctx, c.appID, "SomethingWrong", "post", content)
	if err != nil {
		return nil, err
	}
	if string(resp) == "null" {
		return nil, nil
	}
	var out Message
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func _Service_SomethingWrong_Handler(srv Service) dapr.InvocationHandlerFunc {
	return func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		out = &common.Content{
			ContentType: "application/json",
		}
		resp, methodErr := srv.SomethingWrong(ctx)
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

func RegisterService(s common.Service, srv Service) {
	s.AddServiceInvocationHandler("Echo", errors.ServiceErrorHandler(_Service_Echo_Handler(srv)))
	s.AddServiceInvocationHandler("Hello", errors.ServiceErrorHandler(_Service_Hello_Handler(srv)))
	s.AddServiceInvocationHandler("SomethingWrong", errors.ServiceErrorHandler(_Service_SomethingWrong_Handler(srv)))
}

func NewServiceServer(address string, srv Service) (common.Service, error) {
	s, err := grpc.NewService(address)
	if err != nil {
		return nil, err
	}
	return s, nil
}
