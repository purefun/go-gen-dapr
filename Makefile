.PHONY: echo-service echo-subs pubsub test

echo-service:
	go run ./cmd/go-gen-dapr/main.go -pkg github.com/purefun/go-gen-dapr/examples/echo/service Service

echo-subscriber:
	go run ./cmd/go-gen-dapr/main.go -pkg github.com/purefun/go-gen-dapr/examples/echo/service -target subscriber Subscriptions

pubsub:
	go run ./cmd/go-gen-dapr/main.go -pkg github.com/purefun/go-gen-dapr/examples/pubsub -target pubsub

test:
	go test ./generator/. -count=1 -v
