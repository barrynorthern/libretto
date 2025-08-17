package publisher

import (
	"context"
	"fmt"
)

// PubSubPublisher is a placeholder for a real Pub/Sub implementation.
// For now, it just logs distinctively to differentiate from NOP.
type PubSubPublisher struct{}

func (PubSubPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	fmt.Printf("[pubsub] publish to %s: %d bytes\n", topic, len(data))
	return nil
}

