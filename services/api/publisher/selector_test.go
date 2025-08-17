package publisher

import (
	"context"
	"os"
	"testing"
)

type capturePublisher struct{ called bool }

func (c *capturePublisher) Publish(ctx context.Context, topic string, data []byte) error {
	c.called = true
	return nil
}

func TestSelectDefaultsToNop(t *testing.T) {
	t.Setenv("PUBSUB_ENABLED", "")
	p := Select()
	if _, ok := p.(NopPublisher); !ok {
		t.Fatalf("expected NopPublisher by default")
	}
}

func TestSelectPubSubWhenEnabled(t *testing.T) {
	t.Setenv("PUBSUB_ENABLED", "true")
	p := Select()
	if _, ok := p.(PubSubPublisher); !ok {
		t.Fatalf("expected PubSubPublisher when PUBSUB_ENABLED=true")
	}
}

func TestSmoke(t *testing.T) {
	p := &capturePublisher{}
	if err := Smoke(context.Background(), p); err != nil {
		t.Fatalf("smoke err: %v", err)
	}
	if !p.called {
		t.Fatalf("expected publish to be called")
	}
}

