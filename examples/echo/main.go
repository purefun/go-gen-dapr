package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pkg/errors"
	"github.com/purefun/go-gen-dapr/examples/echo/service"
)

// dapr run -a echo_server -p 6000 -P grpc -- go run . -server
// dapr run -a client -P grpc -- go run . -client

func main() {
	runClient := flag.Bool("client", false, "run client")
	runServer := flag.Bool("server", false, "run server")
	flag.Parse()

	if !*runClient && !*runServer {
		panic("please add --client or --server flag to run the demo")
	}

	if *runClient {
		NewClient()
	}
	if *runServer {
		NewServer()
	}
}

func NewClient() {
	echo, _ := service.NewServiceClient("echo_server")
	resp1, err := echo.Echo(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("echo.Echo returns:", *resp1)

	resp2, err := echo.Hello(context.Background(), service.Message{Text: "Hello"})
	fmt.Println("echo.Hello returns:", resp2.Text)

	_, err = echo.SomethingWrong(context.Background())
	fmt.Println(err)
}

func NewServer() {
	h := new(Handlers)
	s, err := service.NewServiceServer(":6000")
	service.RegisterService(s, h)
	if err != nil {
		panic(err)
	}
	s.Start()
}

type Handlers struct {
}

func (h *Handlers) Hello(ctx context.Context, in service.Message) (*service.Message, error) {
	return &service.Message{Text: in.Text + " world!"}, nil
}

func (h *Handlers) Echo(ctx context.Context) (*string, error) {
	out := "Echo called"
	return &out, nil
}

func (h *Handlers) SomethingWrong(ctx context.Context) (*service.Message, error) {
	return nil, errors.Errorf("something wrong error")
}
