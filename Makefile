.PHONY: echo-service echo-subs

echo-service:
	go run ./cmd/go-gen-dapr/main.go -pkg github.com/purefun/go-gen-dapr/examples/echo/service Service

echo-subs:
	go run ./cmd/go-gen-dapr/main.go -pkg github.com/purefun/go-gen-dapr/examples/echo/service -target subscriptions Subscriptions

pubsub:
	go run ./cmd/go-gen-dapr/main.go -pkg github.com/purefun/go-gen-dapr/examples/pubsub -target pubsub
