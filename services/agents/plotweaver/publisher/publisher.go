package publisher

import (
	"context"
	"fmt"
)

// Publisher publishes events (dev or real). For now, just logs distinctively per implementation.
type Publisher interface {
	Publish(ctx context.Context, topic string, data []byte) error
}

type NopPublisher struct{}

func (NopPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	fmt.Printf("publish to %s: %d bytes\n", topic, len(data))
	return nil
}

type PubSubPublisher struct{}

func (PubSubPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	fmt.Printf("[pubsub] publish to %s: %d bytes\n", topic, len(data))
	return nil
}

type DevPushPublisher struct{}

func (DevPushPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	fmt.Printf("[devpush] publish to %s: %d bytes\n", topic, len(data))
	return nil
}

