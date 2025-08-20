package publisher

import (
	"context"
	"fmt"
	"os"
)

// Select returns a Publisher based on env selection.
// Priority:
// 1) PUBLISHER in {nop, devpush, pubsub}
// 2) Back-compat: if PUBSUB_ENABLED=true -> pubsub
// 3) Default: nop
func Select() Publisher {
	// Primary selector
	if v := os.Getenv("PUBLISHER"); v != "" {
		switch v {
		case "pubsub":
			return PubSubPublisher{}
		case "devpush":
			return newDevPushFromEnv()
		case "nop":
			fallthrough
		default:
			return NopPublisher{}
		}
	}
	// Back-compat path
	if os.Getenv("PUBSUB_ENABLED") == "true" {
		return PubSubPublisher{}
	}
	return NopPublisher{}
}

func newDevPushFromEnv() Publisher {
	url := os.Getenv("PLOT_WEAVER_URL")
	if url == "" {
		port := os.Getenv("PLOT_PORT")
		if port == "" {
			port = "8081"
		}
		url = fmt.Sprintf("http://localhost:%s/push", port)
	}
	return DevPushPublisher{URL: url}
}

// Smoke publishes a small payload to verify the publisher wiring.
func Smoke(ctx context.Context, p Publisher) error {
	return p.Publish(ctx, "libretto.dev.smoke", []byte("ok"))
}
