package publisher

import (
	"context"
	"os"
)

// Select returns a Publisher based on env: if PUBSUB_ENABLED=true, returns PubSubPublisher; else NopPublisher.
func Select() Publisher {
	if os.Getenv("PUBSUB_ENABLED") == "true" {
		return PubSubPublisher{}
	}
	return NopPublisher{}
}

// Smoke publishes a small payload to verify the publisher wiring.
func Smoke(ctx context.Context, p Publisher) error {
	return p.Publish(ctx, "libretto.dev.smoke", []byte("ok"))
}

