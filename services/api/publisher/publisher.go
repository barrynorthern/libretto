package publisher

import (
	"context"
	"fmt"
)

// Publisher publishes events to a bus. In MVP this will be Pub/Sub.
type Publisher interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

// NopPublisher is used in local tests.
type NopPublisher struct{}

func (NopPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	fmt.Printf("publish to %s: %d bytes\n", topic, len(data))
	return nil
}

