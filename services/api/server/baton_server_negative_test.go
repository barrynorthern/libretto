package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	batonv1 "github.com/barrynorthern/libretto/gen/go/libretto/baton/v1"
)

type noopPublisher struct{}

func (noopPublisher) Publish(ctx context.Context, topic string, data []byte) error { return nil }

func TestIssueDirectiveRejectsInvalidEnvelopeVersion(t *testing.T) {
	// Force invalid semver to trigger validation error
	t.Setenv("EVENT_VERSION", "badversion")
	svc := &BatonServer{Pub: noopPublisher{}, Topic: "t", Producer: "api"}
	req := connect.NewRequest(&batonv1.IssueDirectiveRequest{Text: "x"})
	_, err := svc.IssueDirective(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error for invalid event version")
	}
	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", connect.CodeOf(err))
	}
}
