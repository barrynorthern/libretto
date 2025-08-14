package server

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	graphv1 "github.com/barrynorthern/libretto/gen/go/libretto/graph/v1"
)

type fakeStore struct{ version string; count int32; err error }

func (f fakeStore) Apply(parent string, deltas []*graphv1.Delta) (string, int32, error) {
	return f.version, f.count, f.err
}

func TestApplySuccess(t *testing.T) {
	s := &GraphWriteServer{Store: fakeStore{version: "01JF00", count: 2}}
	req := connect.NewRequest(&graphv1.ApplyRequest{ParentVersionId: "01JROOT", Deltas: []*graphv1.Delta{{Op: "create"}, {Op: "create"}}})
	res, err := s.Apply(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got, want := res.Msg.GetGraphVersionId(), "01JF00"; got != want {
		t.Fatalf("version got %q want %q", got, want)
	}
	if got, want := res.Msg.GetApplied(), int32(2); got != want {
		t.Fatalf("applied got %d want %d", got, want)
	}
}

