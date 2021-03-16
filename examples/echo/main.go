package main

import (
	"context"
	"flag"
	"fmt"
)

// dapr run -H 3500 -a echo_server -p 6000 -P grpc -- go run . --server
// dapr run -a client -p 6000 -P grpc -- go run . -client

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
	echo, _ := NewServiceClient("echo_server")
	resp1, err := echo.Echo(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("echo.Echo returns:", *resp1)

	resp2, err := echo.Hello(context.Background(), Message{Text: "Hello"})
	fmt.Println("echo.Hello returns:", resp2.Text)
}

func NewServer() {
	h := new(Handlers)
	s, err := NewServiceServer(":6000", h)
	if err != nil {
		panic(err)
	}
	s.Start()
}

type Handlers struct {
}

func (h *Handlers) Hello(ctx context.Context, in Message) (*Message, error) {
	return &Message{Text: in.Text + " world!"}, nil
}

func (h *Handlers) Echo(ctx context.Context) (*string, error) {
	out := "Echo called"
	return &out, nil
}
